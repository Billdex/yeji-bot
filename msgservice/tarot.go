package msgservice

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/dao"
	"yeji-bot/pkg/kit"
)

func QueryTarot(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) error {
	y, m, d := time.Now().Date()
	timeSeed := time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix()
	authorMd5 := md5.Sum([]byte(msg.Author.MemberOpenid))
	seed := timeSeed + int64(binary.BigEndian.Uint64(authorMd5[:]))
	selfRand := rand.New(rand.NewSource(seed))

	tarots, err := dao.ListAllTarots(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("list all tarots fail. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "获取签文信息失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}
	tarot := tarots[selfRand.Intn(len(tarots))]
	if tarot.Score == 99 && msg.Author.MemberOpenid != "E4F38230B10C96011C4B7A7D87B008E8" {
		score := selfRand.Int63n(98)
		for i := range tarots {
			if int64(tarots[i].Score) == score {
				tarot = tarots[i]
				break
			}
		}
	}

	content := fmt.Sprintf("签筒缓缓落出一支签:\n[%d %s] %s", tarot.Score, tarot.LevelStr(), tarot.Description)
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: content,
		MsgType: openapi.MsgTypeText,
		MsgId:   msg.Id,
		MsgSeq:  kit.Seq(ctx),
	})
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}
