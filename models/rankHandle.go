package models

import (
    "healing2020/models/statements"
    "healing2020/pkg/setting"

    //"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    //"fmt"
    //"time"
    "encoding/json"
    "strconv"
)

type Rank struct{
    User string `json:"user"`
    Avatar string `json:"avatar"`

    Type int `json:"type"`
    Photo string `json:"photo"`
    Text string `json:"text"`
    Source string `json:"source"`

    Time string `json:"time"`
    Praise int `json:"praise"`
    Name string `json:"name"`
}

type AllRank struct {
    Time string `json:"time"`
    Data []Rank `json:"data"`
}

func SendDeliverRank(){
    //everyday's work
    
    //get from mysql
    db := setting.MysqlConn()
    defer db.Close()
    var deliver []statements.Deliver
    result := db.Model(&statements.Deliver{}).Where("created_at LIKE ?","2021_02_09%").Order("created_at desc, praise").Find(&deliver)
    rows := result.RowsAffected
    var rank []Rank = make([]Rank,10)
    for i:=0;i<int(rows);i++ {
        rank[i].Type = deliver[i].Type
        rank[i].Text = deliver[i].TextField
        rank[i].Photo = deliver[i].Photo
        rank[i].Source = deliver[i].Record

        userid := deliver[i].UserId
        var user statements.User
        db.Model(&statements.User{}).Where("id=?",userid).First(&user)
        rank[i].User = user.NickName
        rank[i].Avatar = user.Avatar
    }
    jsonRank,_ := json.Marshal(rank)

    //set in redis
    client := setting.RedisConn()
    defer client.Close()
    count,_ := client.Get("rankCount").Float64()
    keyName := "Deliver." + strconv.FormatFloat(count,'f',2,64)
    client.Set(keyName,jsonRank,0)
    count = count + 1
    client.Set("rankCount",count,0)
    redisRank,_ := client.Get("2.22").Bytes()
    var rank2 []Rank
    json.Unmarshal(redisRank,&rank2)
    //fmt.Println(rank2)
}

func GetDeliverRank() ([]AllRank,string){
    var result []AllRank
    client := setting.RedisConn()
    defer client.Close()
    count,_ := client.Get("rankCount").Float64()
    var i float64 = 0
    for j:=0;;j++ {
        var rank []Rank
        if i>=count {
            break
        }
        var date float64 = 3.15+i
        dateStr := strconv.FormatFloat(date,'f',2,64)
        dateStr = "User." + dateStr
        data,err := client.Get(dateStr).Bytes()
        if err != nil {
            return nil,"Unexpected data" 
        }
        json.Unmarshal(data,&rank)
        i=i+0.01

        result[j].Data = rank
    }

    //fmt.Println(result)
    return result,""
}

func GetSongRank() ([]AllRank,string){
    var result []AllRank
    client := setting.RedisConn()
    defer client.Close()
    count,_ := client.Get("rankCount").Float64()
    var i float64 = 0
    for j:=0;;j++ {
        var rank []Rank
        if i>=count {
            break
        }
        var date float64 = 3.15+i
        dateStr := strconv.FormatFloat(date,'f',2,64)
        dateStr = "Song." + dateStr
        data,err := client.Get(dateStr).Bytes()
        if err != nil {
            return nil,"Unexpected data" 
        }
        json.Unmarshal(data,&rank)
        i=i+0.01

        result[j].Data = rank
    }

    //fmt.Println(result)
    return result,""
}

func GetAllUserRank() ([]AllRank,string){
    var result []AllRank
    client := setting.RedisConn()
    defer client.Close()
    count,_ := client.Get("rankCount").Float64()
    var i float64 = 0
    for j:=0;;j++ {
        var rank []Rank
        if i>=count {
            break
        }
        var date float64 = 3.15+i
        dateStr := strconv.FormatFloat(date,'f',2,64)
        dateStr = "Song." + dateStr
        data,err := client.Get(dateStr).Bytes()
        if err != nil {
            return nil,"Unexpected data" 
        }
        json.Unmarshal(data,&rank)
        i=i+0.01

        result[j].Data = rank
    }

    //fmt.Println(result)
    return result,""
}

type UserRank struct {
    Rank int `json:"rank"`
}

func GetUserRank(id string) (UserRank,error){
    intId,_ := strconv.Atoi(id)
    userId := uint(intId)
    db := setting.MysqlConn()
    defer db.Close()

    rows,err := db.Model(&statements.User{}).Order("Money desc").Rows() 
    rank := 0
    if err != nil {
        return UserRank{},err
    }
    for rows.Next() {
        var user statements.User
        db.ScanRows(rows,&user)
        if user.ID == userId {
            break
        }
        rank++
    }
    var userRank UserRank
    userRank.Rank = rank
    return userRank,err
}
