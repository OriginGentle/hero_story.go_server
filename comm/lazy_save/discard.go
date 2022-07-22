package lazy_save

import "hero_story.go_server/comm/log"

func Discard(lso LazySaveObj) {
	if nil == lso {
		return
	}

	log.Info("放弃延时保存，lsoId = %s", lso.GetLsoId())

	lsoMap.Delete(lso.GetLsoId())
}
