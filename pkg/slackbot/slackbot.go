/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package slackbot

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/n3wscott/slacker/pkg/slackbot/binding"
	"github.com/slack-go/slack"
)

type Instance interface {
	GetTeam(ctx context.Context) (*SlackTeamInfo, error)
	GetChannels(ctx context.Context) (SlackChannels, error)
	GetIMs(ctx context.Context) (SlackIMs, error)
	GetUsers(ctx context.Context) (SlackUsers, error)
	PostMessage(ctx context.Context, channelID, message string) (*SlackPostResponse, error)
	PostResponse(ctx context.Context, channelID, threadID, message string) (*SlackPostResponse, error)
	PostReaction(ctx context.Context, channelID, threadID, reaction string) (*SlackPostResponse, error)
	RemoveReaction(ctx context.Context, channelID, threadID, reaction string) (*SlackPostResponse, error)
}

func NewInstance(ctx context.Context) Instance {
	return &slackbotInstance{}
}

type slackbotInstance struct {
	mtx sync.Mutex

	team     *SlackTeamInfo
	channels SlackChannels
	ims      SlackIMs
	users    SlackUsers

	client  sync.Once
	_client *slack.Client
}

var _ Instance = (*slackbotInstance)(nil)

// Start is a blocking call.
func (s *slackbotInstance) getClient(ctx context.Context) *slack.Client {
	s.client.Do(func() {
		c, err := binding.New(ctx)
		if err != nil {
			panic(err)
		}
		s._client = c
	})
	return s._client
}

func (s *slackbotInstance) syncTeam(ctx context.Context) error {
	ti, err := s.getClient(ctx).GetTeamInfo()
	if err != nil {
		return err
	}

	u, err := url.Parse(fmt.Sprintf("https://%s.slack.com/", ti.Domain))
	if err != nil {
		return err
	}

	s.mtx.Lock()
	s.team = &SlackTeamInfo{
		ID:   ti.ID,
		Name: ti.Name,
		URL:  u.String(),
	}
	s.mtx.Unlock()

	return nil
}

func (s *slackbotInstance) GetTeam(ctx context.Context) (*SlackTeamInfo, error) {
	if s.team == nil {
		if err := s.syncTeam(ctx); err != nil {
			return nil, err
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.team == nil {
		return nil, errors.New("Team not set.")
	}
	return s.team, nil
}

func (s *slackbotInstance) syncChannels(ctx context.Context) error {
	chs, err := s.getClient(ctx).GetChannels(true)
	if err != nil {
		return err
	}

	channels := make(SlackChannels, len(chs))
	for i, ch := range chs {
		channels[i] = SlackChannel{
			Name:     ch.Name,
			ID:       ch.ID,
			IsMember: ch.IsMember,
		}
	}

	s.mtx.Lock()
	s.channels = channels
	s.mtx.Unlock()

	return nil
}

func (s *slackbotInstance) GetChannels(ctx context.Context) (SlackChannels, error) {
	if s.channels == nil {
		if err := s.syncChannels(ctx); err != nil {
			return nil, err
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.channels == nil {
		return nil, errors.New("Channels not set.")
	}
	return s.channels, nil
}

func (s *slackbotInstance) syncIMs(ctx context.Context) error {
	users, err := s.GetUsers(ctx)
	if err != nil {
		return err
	}

	userName := make(map[string]string, len(users))
	for _, u := range users {
		userName[u.ID] = u.Name
	}

	imcs, err := s.getClient(ctx).GetIMChannels()
	if err != nil {
		return err
	}

	ims := make(SlackIMs, 0, len(imcs))
	for _, imc := range imcs {
		if imc.IsUserDeleted {
			continue
		}

		ims = append(ims, SlackIM{
			Name: userName[imc.User],
			ID:   imc.ID,
			With: imc.User,
		})
	}

	s.mtx.Lock()
	s.ims = ims
	s.mtx.Unlock()

	return nil
}

func (s *slackbotInstance) syncUsers(ctx context.Context) error {
	us, err := s.getClient(ctx).GetUsers()
	if err != nil {
		return err
	}

	users := make(SlackUsers, 0, len(us))
	for _, u := range us {
		users = append(users, SlackUser{
			Name: u.Name,
			ID:   u.ID,
		})
	}

	s.mtx.Lock()
	s.users = users
	s.mtx.Unlock()

	return nil
}

func (s *slackbotInstance) GetUsers(ctx context.Context) (SlackUsers, error) {
	if s.users == nil {
		if err := s.syncUsers(ctx); err != nil {
			return nil, err
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.users == nil {
		return nil, errors.New("Users not set.")
	}
	return s.users, nil
}

func (s *slackbotInstance) GetIMs(ctx context.Context) (SlackIMs, error) {
	if s.ims == nil {
		if err := s.syncIMs(ctx); err != nil {
			return nil, err
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.ims == nil {
		return nil, errors.New("IMs not set.")
	}
	return s.ims, nil
}

func (s *slackbotInstance) PostMessage(ctx context.Context, channelID, message string) (*SlackPostResponse, error) {
	channelID, timestamp, err := s.getClient(ctx).PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	return &SlackPostResponse{
		ChannelID: channelID,
		Thread:    timestamp,
		Parent:    timestamp,
	}, nil
}

func (s *slackbotInstance) PostResponse(ctx context.Context, channelID, parentID, message string) (*SlackPostResponse, error) {
	client := s.getClient(ctx)

	channelID, threadID, err := client.PostMessage(channelID, slack.MsgOptionText(message, false), slack.MsgOptionTS(parentID))
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	return &SlackPostResponse{
		ChannelID: channelID,
		Thread:    threadID,
		Parent:    parentID,
	}, nil
}

func (s *slackbotInstance) PostReaction(ctx context.Context, channelID, threadID, reaction string) (*SlackPostResponse, error) {
	client := s.getClient(ctx)

	msgRef := slack.NewRefToMessage(channelID, threadID)

	err := client.AddReaction(reaction, msgRef)
	if err != nil {
		return nil, err
	}

	return &SlackPostResponse{
		ChannelID: channelID,
		Thread:    threadID,
	}, nil
}

func (s *slackbotInstance) RemoveReaction(ctx context.Context, channelID, threadID, reaction string) (*SlackPostResponse, error) {
	client := s.getClient(ctx)

	msgRef := slack.NewRefToMessage(channelID, threadID)

	err := client.RemoveReaction(reaction, msgRef)
	if err != nil {
		return nil, err
	}

	return &SlackPostResponse{
		ChannelID: channelID,
		Thread:    threadID,
	}, nil
}
