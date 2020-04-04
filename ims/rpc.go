/**
 * Copyright (c) 2014-2015, GoBelieve
 * All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

package ims

import (
	"github.com/GoBelieveIO/im_service/core"
	"sync/atomic"
)

func SyncMessage(addr string, sync_key *core.SyncHistory) *core.PeerHistoryMessage {
	atomic.AddInt64(&server_summary.nrequests, 1)
	messages, last_msgid, hasMore := storage.LoadHistoryMessagesV3(sync_key.AppID, sync_key.Uid, sync_key.LastMsgID, config.limit, config.hard_limit)

	historyMessages := make([]*core.HistoryMessage, 0, 10)
	for _, emsg := range messages {
		hm := &core.HistoryMessage{}
		hm.MsgID = emsg.MsgId
		hm.DeviceID = emsg.DeviceId
		hm.Cmd = int32(emsg.Msg.Cmd)

		emsg.Msg.Version = core.DEFAULT_VERSION
		hm.Raw = emsg.Msg.ToData()
		historyMessages = append(historyMessages, hm)
	}

	return &core.PeerHistoryMessage{historyMessages, last_msgid, hasMore}
}

func SyncGroupMessage(addr string, sync_key *core.SyncGroupHistory) *core.GroupHistoryMessage {
	atomic.AddInt64(&server_summary.nrequests, 1)
	messages, last_msgid := storage.LoadGroupHistoryMessages(sync_key.AppID, sync_key.Uid, sync_key.GroupID, sync_key.LastMsgID, sync_key.Timestamp, GROUP_OFFLINE_LIMIT)

	historyMessages := make([]*core.HistoryMessage, 0, 10)
	for _, emsg := range messages {
		hm := &core.HistoryMessage{}
		hm.MsgID = emsg.MsgId
		hm.DeviceID = emsg.DeviceId
		hm.Cmd = int32(emsg.Msg.Cmd)

		emsg.Msg.Version = core.DEFAULT_VERSION
		hm.Raw = emsg.Msg.ToData()
		historyMessages = append(historyMessages, hm)
	}

	return &core.GroupHistoryMessage{historyMessages, last_msgid, false}
}

func SavePeerMessage(addr string, m *core.PeerMessage) ([2]int64, error) {
	atomic.AddInt64(&server_summary.nrequests, 1)
	atomic.AddInt64(&server_summary.peer_message_count, 1)
	msg := &core.Message{Cmd: int(m.Cmd), Version: core.DEFAULT_VERSION}
	msg.FromData(m.Raw)
	msgid, prev_msgid := storage.SavePeerMessage(m.AppID, m.Uid, m.DeviceID, msg)
	return [2]int64{msgid, prev_msgid}, nil
}

func SavePeerGroupMessage(addr string, m *core.PeerGroupMessage) ([]int64, error) {
	atomic.AddInt64(&server_summary.nrequests, 1)
	atomic.AddInt64(&server_summary.peer_message_count, 1)
	msg := &core.Message{Cmd: int(m.Cmd), Version: core.DEFAULT_VERSION}
	msg.FromData(m.Raw)
	r := storage.SavePeerGroupMessage(m.AppID, m.Members, m.DeviceID, msg)
	return r, nil
}

func SaveGroupMessage(addr string, m *core.GroupMessage) ([2]int64, error) {
	atomic.AddInt64(&server_summary.nrequests, 1)
	atomic.AddInt64(&server_summary.group_message_count, 1)
	msg := &core.Message{Cmd: int(m.Cmd), Version: core.DEFAULT_VERSION}
	msg.FromData(m.Raw)
	msgid, prev_msgid := storage.SaveGroupMessage(m.AppID, m.GroupID, m.DeviceID, msg)

	return [2]int64{msgid, prev_msgid}, nil
}

func GetNewCount(addr string, sync_key *core.SyncHistory) (int64, error) {
	atomic.AddInt64(&server_summary.nrequests, 1)
	count := storage.GetNewCount(sync_key.AppID, sync_key.Uid, sync_key.LastMsgID)
	return int64(count), nil
}

func GetLatestMessage(addr string, r *core.HistoryRequest) []*core.HistoryMessage {
	atomic.AddInt64(&server_summary.nrequests, 1)
	messages := storage.LoadLatestMessages(r.AppID, r.Uid, int(r.Limit))

	historyMessages := make([]*core.HistoryMessage, 0, 10)
	for _, emsg := range messages {
		hm := &core.HistoryMessage{}
		hm.MsgID = emsg.MsgId
		hm.DeviceID = emsg.DeviceId
		hm.Cmd = int32(emsg.Msg.Cmd)

		emsg.Msg.Version = core.DEFAULT_VERSION
		hm.Raw = emsg.Msg.ToData()
		historyMessages = append(historyMessages, hm)
	}

	return historyMessages
}
