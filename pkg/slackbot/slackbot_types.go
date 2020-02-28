/*
Copyright 2020 The Knative Authors.

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

type SlackTeamInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"domain,omitempty"`
}

type SlackChannel struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	IsMember bool   `json:"isMember,omitempty"`
}

type SlackChannels []SlackChannel

type SlackIM struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
	With string `json:"with,omitempty"`
}

type SlackIMs []SlackIM

type SlackUser struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type SlackUsers []SlackUser

type SlackPostResponse struct {
	ChannelID string
	Parent    string
	Thread    string
}
