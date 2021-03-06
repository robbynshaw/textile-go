package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
)

var notificationStore repo.NotificationStore

func init() {
	setupNotificationDB()
}

func setupNotificationDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	notificationStore = NewNotificationStore(conn, new(sync.Mutex))
}

func TestNotificationDB_Add(t *testing.T) {
	err := notificationStore.Add(&repo.Notification{
		Id:        "abcde",
		Date:      time.Now(),
		ActorId:   ksuid.New().String(),
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   ksuid.New().String(),
		Type:      repo.InviteReceivedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "abcde" {
		t.Errorf(`expected "abcde" got %s`, id)
	}
}

func TestNotificationDB_Get(t *testing.T) {
	notif := notificationStore.Get("abcde")
	if notif == nil {
		t.Error("could not get notification")
	}
}

func TestNotificationDB_Read(t *testing.T) {
	err := notificationStore.Read("abcde")
	if err != nil {
		t.Error(err)
		return
	}
	notifs := notificationStore.List("", 1)
	if len(notifs) == 0 || !notifs[0].Read {
		t.Error("notification read bad result")
	}
}

func TestNotificationDB_ReadAll(t *testing.T) {
	setupNotificationDB()
	err := notificationStore.Add(&repo.Notification{
		Id:        "abcde",
		Date:      time.Now(),
		ActorId:   ksuid.New().String(),
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   ksuid.New().String(),
		Type:      repo.InviteReceivedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&repo.Notification{
		Id:        "abcdef",
		Date:      time.Now(),
		ActorId:   ksuid.New().String(),
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   ksuid.New().String(),
		Type:      repo.PeerJoinedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.ReadAll()
	if err != nil {
		t.Error(err)
		return
	}
	notifs := notificationStore.List("", -1)
	if len(notifs) != 2 || !notifs[0].Read || !notifs[1].Read {
		t.Error("notification read all bad result")
	}
}

func TestNotificationDB_List(t *testing.T) {
	setupNotificationDB()
	err := notificationStore.Add(&repo.Notification{
		Id:        "abc",
		Date:      time.Now(),
		ActorId:   "actor1",
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   "block1",
		Type:      repo.InviteReceivedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&repo.Notification{
		Id:        "def",
		Date:      time.Now().Add(time.Minute),
		ActorId:   "actor1",
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   "block2",
		Type:      repo.PeerJoinedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&repo.Notification{
		Id:        "ghi",
		Date:      time.Now().Add(time.Minute * 2),
		ActorId:   "actor2",
		Subject:   "test",
		SubjectId: ksuid.New().String(),
		BlockId:   "block2",
		Type:      repo.CommentAddedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&repo.Notification{
		Id:        "jkl",
		Date:      time.Now().Add(time.Minute * 3),
		ActorId:   "actor3",
		Subject:   "test",
		SubjectId: "subject1",
		BlockId:   "block3",
		Target:    "target",
		Type:      repo.FilesAddedNotification,
	})
	if err != nil {
		t.Error(err)
	}
	all := notificationStore.List("", -1)
	if len(all) != 4 {
		t.Error("returned incorrect number of notifications")
		return
	}
	limited := notificationStore.List("", 1)
	if len(limited) != 1 {
		t.Error("returned incorrect number of notifications")
		return
	}
	offset := notificationStore.List(limited[0].Id, -1)
	if len(offset) != 3 {
		t.Error("returned incorrect number of notifications")
		return
	}
}

func TestNotificationDB_CountUnread(t *testing.T) {
	cnt := notificationStore.CountUnread()
	if cnt != 4 {
		t.Error("returned incorrect count of unread notifications")
	}
}

func TestNotificationDB_Delete(t *testing.T) {
	err := notificationStore.Delete("abc")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abc").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteByActor(t *testing.T) {
	err := notificationStore.DeleteByActor("actor1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("def").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteBySubject(t *testing.T) {
	err := notificationStore.DeleteBySubject("subject1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("jkl").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteByBlock(t *testing.T) {
	err := notificationStore.DeleteByBlock("block2")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("ghi").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
