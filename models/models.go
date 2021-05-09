package models

import (
    "healing2020/models/statements"
)

//默认约定，直接嵌入gorm.Model即可
//type Model struct {
//    ID uint `gorm:"primaryKey"`
//    CreatedAt time.Time
//    UpdatedAt time.Time
//    DeletedAt gorm.DeletedAt `gorm:"index"`
//}

func TableInit() {
    statements.UserInit()
    statements.CommentInit()
    statements.DeliverInit()
    statements.LotteryInit()
    statements.MailboxInit()
    statements.MessageInit()
    statements.PrizeInit()
    statements.SongInit()
    statements.SpecialInit()
    statements.SubjectInit()
    statements.VodInit()
    statements.UserOtherInit()
    statements.PraiseInit()
}
