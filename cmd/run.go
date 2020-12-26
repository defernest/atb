/*
Copyright © 2020 Denis Nesterenko <defernest@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/defernest/atb/pkg/gitlab"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := NewATBot()
		if err != nil {
			log.Err(err).Msg("Ошибка при создании экземпляра бота")
		}
		b.SetToken()
		err = b.SetUrl()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		b.Bot.Handle("/start", func(m *tb.Message) {
			b.Bot.Send(m.Sender,"Привет! Я бот для проведения быстрого голосования в концепции SCRUM")
		})
		b.Bot.Handle("/help", func(m *tb.Message) {
			b.Bot.Send(m.Sender, "Я помогу вам быстро оценить задачи по спринту и автоматически проставлю их estimate в Gitlab'е\\.", &tb.SendOptions{ParseMode: tb.ModeMarkdownV2, DisableNotification: true})
			b.Bot.Send(m.Sender, "Для получения списка задач мне можно прислать URL на майлстоун, название майлстоуна или его IID: \n" +
				"`\\/sprint http[s]:\\/\\/gitlab-url.com/Project` \\| `\\/sprint Название майлстоуна` \\| `\\/sprint 123`", &tb.SendOptions{ParseMode: tb.ModeMarkdownV2, DisableNotification: true})
			b.Bot.Send(m.Sender, "Для старта процедуры голосования пришли мне IID задачи, в Gitlab'е: \n" +
				"`/vote 1234`", &tb.SendOptions{ParseMode: tb.ModeMarkdownV2, DisableNotification: true})
			b.Bot.Send(m.Sender, "По завершению процесса голосования нужно остановить опрос\\! Сделать это можно соответствующей кнопкой под опросом`", &tb.SendOptions{ParseMode: tb.ModeMarkdownV2, DisableNotification: true})
		})
		b.Bot.Handle("/sprint", func(m *tb.Message) {
			issues, err := gitlab.GetListIssues(m.Payload, b.Token())
			if err != nil {
				b.Bot.Send(m.Chat, err)
			}
			var issuesResult string
			for _, issue := range issues {
				issuesResult = fmt.Sprintf("%s\n[ %d ] %s (%d)", issuesResult, issue.IID, issue.Title, issue.TimeStats.TimeEstimate/3600)
			}
			b.Bot.Send(m.Chat, issuesResult)
		})
		b.Bot.Handle("/vote", func(m *tb.Message) {
			b.Bot.Send(m.Chat, "голосование")
		})
		b.Bot.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type ATBot struct {
	Bot *tb.Bot
	token string
	url url.URL
}

func (bot *ATBot) Url() url.URL {
	return bot.url
}

func (bot *ATBot) SetUrl() error {
	url, err := parseRawURL(viper.GetString("gitlab-project"))
	if err != nil {
		return err
	}
	bot.url = url
	return nil
}

func (bot *ATBot) Token() string {
	return bot.token
}

func (bot *ATBot)SetToken() {
	bot.token = viper.GetString("gitlab-token")
}

func NewATBot() (*ATBot, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  viper.GetString("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Err(err).Str("service", "bot").Msg("Failed to start Telegram Bot")
	}

	return &ATBot{Bot: bot}, nil
}

func parseRawURL(s string) (url.URL, error) {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return url.URL{}, fmt.Errorf("url is not valid (%s)", s)
	}
	return *u, nil
}