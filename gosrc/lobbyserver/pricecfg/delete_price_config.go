package pricecfg

// DeletePriceCfg 从配置列表中移除价格配置表
func DeletePriceCfg(roomType int, priceCfgType string) {
	if roomType == 0 || priceCfgType == "" {
		return
	}

	cfg, ok := cfgs[roomType]
	if !ok {
		return
	}

	if priceCfgType == "originalPrice" {
		cfg.OriginalPriceCfg = nil
	} else {
		cfg.ActivityPriceCfg = nil
	}

	cfgs[roomType] = cfg

	if cfg.OriginalPriceCfg == nil && cfg.ActivityPriceCfg == nil {
		delete(cfgs, roomType)
		cfg = nil
	}

	if cfg == nil {
		// 只有在原价表与活动表都没有的情况下才删除配置
		notifyAllUserPriceChange(roomType, cfg, deletePriceCfg)
	} else {
		// 原价表与活动表存在任何一种都只是更新配置
		notifyAllUserPriceChange(roomType, cfg, updatePriceCfg)
	}
}

// PushAllPlayer 推送数据给所有用户
// func pushAllPlayer(ops int32, data []byte) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	list, err := redis.Strings(conn.Do("SUNION", stateless.PushServerList))
// 	if err != nil {
// 		log.Println("GetPushServerList err:", err)
// 		return
// 	}

// 	for _, v := range list {
// 		id := uuid.NewV4()
// 		key := fmt.Sprintf("client:%d:%d:%s", 0, ops, id)

// 		conn.Send("MULTI")
// 		conn.Send("SET", key, data)
// 		conn.Send("PUBLISH", v, key)
// 		_, err := conn.Do("EXEC")
// 		if err != nil {
// 			log.Printf("push allplayer server:%s event %s err:%v", v, key, err)
// 			continue
// 		}
// 		log.Printf("push allplayer server:%s push event %s", v, key)
// 	}
// }
