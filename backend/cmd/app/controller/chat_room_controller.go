package controller

import (
	"net/http"

	"github.com/Hibagon1go/ChatApp_Go_React/cmd/app/auth"
	"github.com/Hibagon1go/ChatApp_Go_React/cmd/app/model"
	"github.com/Hibagon1go/ChatApp_Go_React/cmd/app/utils"
	"github.com/labstack/echo/v4"
)

type buildChatRoomReq struct {
	RoomName string `json:"room_name"`
}

type joinNewRoomReq struct {
	ChatRoomID string `json:"chat_room_id"`
}

type updateRoomNameReq struct {
	ChatRoomID string `json:"chat_room_id"`
	RoomName   string `json:"room_name"`
}

type leaveNewRoomReq struct {
	ChatRoomID string `json:"chat_room_id"`
}

func BuildChatRoom(c echo.Context) error {
	chatRoomID := utils.GenerateUUID()
	userID := auth.UserAuth(c)

	req := new(buildChatRoomReq)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	roomName := req.RoomName

	chatRoom := model.ChatRoom{
		ID:   chatRoomID,
		Name: roomName,
	}

	userChatRoom := model.UserChatRoom{
		UserID:     userID,
		ChatRoomID: chatRoomID,
	}

	if err := chatRoom.Create().Error; err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "このチャンネル名は既に使われています"})
	}

	userChatRoom.Create()

	// ルーム作成直後に入れるように、ルームIDを返す
	return c.JSON(http.StatusOK, echo.Map{"chat_room_id": chatRoomID})
}

func JoinNewRoom(c echo.Context) error {
	userID := auth.UserAuth(c)

	req := new(joinNewRoomReq)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	joinNewRoomID := req.ChatRoomID

	userChatRoom := model.UserChatRoom{
		UserID:     userID,
		ChatRoomID: joinNewRoomID,
	}

	if err := userChatRoom.Create().Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "既にこのチャンネルには参加しています"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "ルームの参加に成功しました"})
}

func FetchAllRooms(c echo.Context) error {
	auth.UserAuth(c)

	allRooms := model.ChatRooms{}
	allRooms.FetchAllRooms()

	return c.JSON(http.StatusOK, allRooms)
}

func FetchJoiningRooms(c echo.Context) error {
	userID := auth.UserAuth(c)

	joiningRooms := model.JoiningRooms{}
	joiningRooms.FetchJoiningRooms(userID)

	return c.JSON(http.StatusOK, joiningRooms)
}

func UpdateRoomName(c echo.Context) error {
	auth.UserAuth(c)

	req := new(updateRoomNameReq)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	chatRoomID := req.ChatRoomID
	newRoomName := req.RoomName

	renamedRoom := model.ChatRoom{
		ID:   chatRoomID,
		Name: newRoomName,
	}

	renamedRoom.UpdateName()

	return c.JSON(http.StatusOK, renamedRoom)
}

func LeaveChatRoom(c echo.Context) error {
	userID := auth.UserAuth(c)

	req := new(leaveNewRoomReq)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	chatRoomID := req.ChatRoomID

	leaveRoom := model.UserChatRoom{}

	leaveRoom.LeaveRoom(userID, chatRoomID)

	return c.JSON(http.StatusOK, echo.Map{"message": "ルームを退出しました"})
}
