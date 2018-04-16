package internal

import (
	"server/msg"
	"github.com/golang/glog"
)

func (r *Room) joinRoom(m *msg.JoinRoom, o *Occupant) {
	if o.room != nil {
		for k, v := range r.occupants {
			if v.Uid == o.Uid {
				// todo 掉线重连现场数据替换处理
				o.Replace(r.occupants[k])
				r.occupants[k] = o

				if o != v {
					v.Close()
					glog.Infoln("掉线重连处理")
				} else {
					glog.Infoln("同一个链接重复请求加入房间")
				}

				r.WriteMsg(&msg.UserInfo{Uid: o.Uid}, o.Uid)
				return
			}
		}
	}
	glog.Errorln(o)
	rinfo := &msg.RoomInfo{
		Number: r.Number,
	}
	userinfos := make([]*msg.UserInfo, 0, r.Cap())
	r.Each(0, func(o *Occupant) bool {
		userinfo := &msg.UserInfo{
			Nickname: o.Nickname,
			Uid:      o.Uid,
			Account:  o.Account,
			Sex:      o.Sex,
			Profile:  o.Profile,
			Chips:    o.Chips,
		}
		userinfos = append(userinfos, userinfo)
		return true
	})

	pos := r.addOccupant(o)

	// 坐下失败转为旁观
	if pos == 0 {
		r.addObserve(o)
	} else {
		userInfo := &msg.UserInfo{
			Nickname: o.Nickname,
			Uid:      o.Uid,
			Account:  o.Account,
			Sex:      o.Sex,
			Profile:  o.Profile,
			Chips:    o.Chips,
		}
		r.WriteMsg(&msg.JoinRoomBroadcast{UserInfo: userInfo}, o.Uid)
	}

	o.RoomID = r.Number
	o.UpdateRoomId()
	o.room = r
	o.WriteMsg(&msg.JoinRoomResp{UserInfos: userinfos, RoomInfo: rinfo})

	glog.Errorln("joinRoom", m)
}

func (r *Room) leaveRoom(m *msg.LeaveRoom, o *Occupant) {
	if o.IsGameing() {
		return
	}

	r.removeObserve(o)
	r.removeOccupant(o)
	o.RoomID = ""
	o.room = nil
	o.UpdateRoomId()
	leave := &msg.LeaveRoom{
		RoomNumber: r.Number,
		Uid:        o.Uid,
	}
	r.WriteMsg(leave)
	glog.Errorln("leaveRoom", m)
}

func (r *Room) bet(m *msg.Bet, o *Occupant) {
	if !o.IsGameing() {
		o.WriteMsg(msg.MSG_NOT_NOT_START)
		return
	}
	o.SetAction(m.Value)
	glog.Errorln("bet", m)
}

func (r *Room) sitDown(m *msg.SitDown, o *Occupant) {
	pos := r.addOccupant(o)
	if pos == 0 {
		r.addObserve(o)
	} else {

	}
	r.WriteMsg(&msg.SitDown{Uid: o.Uid, Pos: o.Pos})

	glog.Errorln("sitDown", m)
}

func (r *Room) standUp(m *msg.StandUp, o *Occupant) {
	o.SetAction(-1)
	r.removeOccupant(o)

	r.addObserve(o)
	r.WriteMsg(&msg.StandUp{Uid: o.Uid})

	glog.Errorln("standUp", m)
}

func (r *Room) fold(m *msg.Fold, o *Occupant) {
	if !o.IsGameing() {
		o.WriteMsg(msg.MSG_NOT_NOT_START)
		return
	}


	o.SetAction(-1)
	o.SetSitdown()
	//r.WriteMsg(&msg.StandUp{Uid: o.Uid})
	glog.Errorln("standUp", m)
}
