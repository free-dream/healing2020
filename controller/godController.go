package controller

import (
    "healing2020/models"
)

func autoCreate() {
    models.CreateUsers("","robot","truename","more","avatar","华工","12312312311",1,"摇滚",100,0,0,0)
    models.CreateUsers("","robot2","truename2","more","avatar2","华工","12312312312",1,"古风",30,1,0,0)
    models.CreateUsers("","robot3","truename3","more","avatar3","中大","12182312312",0,"古风",130,1,1,0)
    models.CreateUsers("","robot4","truename4","more","avatar4","中大","18182311310",1,"日语",120,0,1,0)
    models.CreateUsers("","robot5","truename5","more","avatar5","华工","12312123412",0,"爵士",300,1,1,1)

    models.CreateBgs("1",1,1,1,1,1,1)
    models.CreateBgs("2",2,1,1,1,1,1)
    models.CreateBgs("3",3,1,1,1,1,1)
    models.CreateBgs("4",4,1,1,1,1,1)
    models.CreateBgs("5",5,1,1,1,1,1)

    models.CreateVods("1","no more","abc","tobor","","")
    models.CreateVods("2","no more","abcd","2tobor","","")
    models.CreateVods("3","more","abde","3tobor","","")
    models.CreateVods("4","more","bde","4tobor","","")
    models.CreateVods("5","more","bdecc","3tobor","","")

    models.CreateSongs("1","2","2","abcd",0,"source1","","")
    models.CreateSongs("3","2","2","abcd",0,"source2","","")
    models.CreateSongs("3","1","1","abc",0,"source3","","")
    models.CreateSongs("4","1","1","abc",0,"source4","","")
    models.CreateSongs("5","3","3","abde",4,"source5","","")
    models.CreateSongs("2","5","5","bdecc",3,"source6","","")

    models.CreatePraises("1",1,"5")
    models.CreatePraises("2",1,"5")
    models.CreatePraises("4",1,"5")
    models.CreatePraises("3",1,"5")
    models.CreatePraises("3",1,"6")
    models.CreatePraises("1",1,"6")
    models.CreatePraises("2",1,"6")
    models.CreatePraises("4",2,"1")

    models.CreateDelivers("2",1,"i am robot","","",1)
    models.CreateDelivers("3",2,"i am robot2","photo1","",0)
    models.CreateDelivers("3",3,"i am robot2","","source7",0)

    models.CreateComments("5",1,"3","","wonderful")
    models.CreateComments("5",2,"","1","me too")
}

func LoadTestData() {
    models.TableCleanUp()
    autoCreate()
}