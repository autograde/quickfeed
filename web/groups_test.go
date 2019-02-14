package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/autograde/aguis/models"
	"github.com/autograde/aguis/scm"
	"github.com/autograde/aguis/web"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func TestDeleteGroup(t *testing.T) {
	const route = "/groups/:gid"

	db, cleanup := setup(t)
	defer cleanup()

	testCourse := models.Course{
		Name:        "Distributed Systems",
		Code:        "DAT520",
		Year:        2018,
		Tag:         "Spring",
		Provider:    "fake",
		DirectoryID: 1,
	}
	admin := createFakeUser(t, db, 1)
	if err := db.CreateCourse(admin.ID, &testCourse); err != nil {
		t.Fatal(err)
	}

	// create user and enroll as student
	user := createFakeUser(t, db, 2)
	if err := db.CreateEnrollment(&models.Enrollment{UserID: user.ID, CourseID: testCourse.ID}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user.ID, testCourse.ID); err != nil {
		t.Fatal(err)
	}

	group := models.Group{CourseID: testCourse.ID}
	if err := db.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	router := echo.NewRouter(e)

	// Add the route to handler.
	router.Add(http.MethodDelete, route, web.DeleteGroup(db))

	requestURL := "/groups/" + strconv.FormatUint(group.ID, 10)
	r := httptest.NewRequest(http.MethodDelete, requestURL, nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)
	// Prepare context with course request.
	router.Find(http.MethodDelete, requestURL, c)

	// Invoke the prepared handler.
	if err := c.Handler()(c); err != nil {
		t.Fatal(err)
	}

	assertCode(t, w.Code, http.StatusOK)
}

func TestGetGroup(t *testing.T) {
	const route = "/groups/:gid"

	db, cleanup := setup(t)
	defer cleanup()

	testCourse := models.Course{
		Name:        "Distributed Systems",
		Code:        "DAT520",
		Year:        2018,
		Tag:         "Spring",
		Provider:    "fake",
		DirectoryID: 1,
	}
	admin := createFakeUser(t, db, 1)
	if err := db.CreateCourse(admin.ID, &testCourse); err != nil {
		t.Fatal(err)
	}

	// create user and enroll as student
	user := createFakeUser(t, db, 2)
	if err := db.CreateEnrollment(&models.Enrollment{UserID: user.ID, CourseID: testCourse.ID}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user.ID, testCourse.ID); err != nil {
		t.Fatal(err)
	}

	group := models.Group{CourseID: testCourse.ID}
	if err := db.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	router := echo.NewRouter(e)

	// Add the route to handler.
	router.Add(http.MethodDelete, route, web.GetGroup(db))

	requestURL := "/groups/" + strconv.FormatUint(group.ID, 10)
	r := httptest.NewRequest(http.MethodGet, requestURL, nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)
	// Prepare context with course request.
	router.Find(http.MethodDelete, requestURL, c)

	// Invoke the prepared handler.
	if err := c.Handler()(c); err != nil {
		t.Fatal(err)
	}
	assertCode(t, w.Code, http.StatusOK)

	var respGroup models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &respGroup); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(respGroup, group) {
		t.Errorf("have response group %+v, while database has %+v", &respGroup, group)
	}
}

func TestPatchGroupStatus(t *testing.T) {
	const route = "/groups/:gid"

	db, cleanup := setup(t)
	defer cleanup()

	course := models.Course{
		Name:        "Distributed Systems",
		Code:        "DAT520",
		Year:        2018,
		Tag:         "Spring",
		Provider:    "fake",
		DirectoryID: 1,
		ID:          1,
	}

	admin := createFakeUser(t, db, 1)
	err := db.CreateCourse(admin.ID, &course)
	if err != nil {
		t.Fatal(err)
	}

	user1 := createFakeUser(t, db, 2)
	user2 := createFakeUser(t, db, 3)

	// enroll users in course and group
	if err := db.CreateEnrollment(&models.Enrollment{
		UserID: user1.ID, CourseID: course.ID, GroupID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user1.ID, course.ID); err != nil {
		t.Fatal(err)
	}
	if err := db.CreateEnrollment(&models.Enrollment{
		UserID: user2.ID, CourseID: course.ID, GroupID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user2.ID, course.ID); err != nil {
		t.Fatal(err)
	}

	group := &models.Group{
		ID:       1,
		CourseID: course.ID,
		Users:    []*models.User{user1, user2},
	}
	err = db.CreateGroup(group)
	if err != nil {
		t.Fatal(err)
	}
	// get the group as stored in db with enrollments
	prePatchGroup, err := db.GetGroup(group.ID)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	router := echo.NewRouter(e)

	// add the route to handler.
	router.Add(http.MethodPatch, route, web.PatchGroup(nullLogger(), db))

	// send empty request, the user should not be modified.
	emptyJSON, err := json.Marshal(&web.UpdateGroupRequest{})
	if err != nil {
		t.Fatal(err)
	}
	requestBody := bytes.NewReader(emptyJSON)

	requestURL := "/groups/" + strconv.FormatUint(group.ID, 10)
	r := httptest.NewRequest(http.MethodPatch, requestURL, requestBody)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)
	f := scm.NewFakeSCMClient()
	if _, err := f.CreateDirectory(context.Background(), &scm.CreateDirectoryOptions{
		Name: course.Code,
		Path: course.Code,
	}); err != nil {
		t.Fatal(err)
	}
	c.Set("fake", f)
	// set admin as the user for this context
	c.Set("user", admin)
	router.Find(http.MethodPatch, requestURL, c)

	// invoke the prepared handler
	if err := c.Handler()(c); err != nil {
		t.Error(err)
	}
	assertCode(t, w.Code, http.StatusOK)

	// check that the group didn't change
	haveGroup, err := db.GetGroup(group.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(prePatchGroup, haveGroup) {
		t.Errorf("have group %+v want %+v", haveGroup, prePatchGroup)
	}

	// send request for status change of the group
	trueJSON, err := json.Marshal(&web.UpdateGroupRequest{Status: 3})
	if err != nil {
		t.Fatal(err)
	}
	requestBody.Reset(trueJSON)

	r = httptest.NewRequest(http.MethodPatch, requestURL, requestBody)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w = httptest.NewRecorder()
	c.Reset(r, w)
	// set admin as the user for this context
	c.Set("user", admin)
	fakeProvider, err := scm.NewSCMClient("fake", "token")
	if err != nil {
		t.Fatal(err)
	}
	fakeProvider.CreateDirectory(c.Request().Context(),
		&scm.CreateDirectoryOptions{Path: "path", Name: "name"},
	)
	c.Set("fake", fakeProvider)
	router.Find(http.MethodPatch, requestURL, c)

	// invoke the prepared handler
	if err := c.Handler()(c); err != nil {
		t.Error(err)
	}
	assertCode(t, w.Code, http.StatusOK)

	// check that the group have changed status
	haveGroup, err = db.GetGroup(group.ID)
	if err != nil {
		t.Fatal(err)
	}
	wantGroup := prePatchGroup
	wantGroup.Status = 3
	if !reflect.DeepEqual(wantGroup, haveGroup) {
		t.Errorf("have group %+v want %+v", haveGroup, wantGroup)
	}
}

func TestGetGroupByUserAndCourse(t *testing.T) {
	db, cleanup := setup(t)
	defer cleanup()

	course := models.Course{
		Name:        "Distributed Systems",
		Code:        "DAT520",
		Year:        2018,
		Tag:         "Spring",
		Provider:    "fake",
		DirectoryID: 1,
		ID:          1,
	}

	admin := createFakeUser(t, db, 1)
	err := db.CreateCourse(admin.ID, &course)
	if err != nil {
		t.Fatal(err)
	}

	user1 := createFakeUser(t, db, 2)
	user2 := createFakeUser(t, db, 3)

	// enroll users in course and group
	if err := db.CreateEnrollment(&models.Enrollment{
		UserID: user1.ID, CourseID: course.ID, GroupID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user1.ID, course.ID); err != nil {
		t.Fatal(err)
	}
	if err := db.CreateEnrollment(&models.Enrollment{
		UserID: user2.ID, CourseID: course.ID, GroupID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.EnrollStudent(user2.ID, course.ID); err != nil {
		t.Fatal(err)
	}

	group := &models.Group{
		ID:       1,
		CourseID: course.ID,
		Users:    []*models.User{user1, user2},
	}
	err = db.CreateGroup(group)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	router := echo.NewRouter(e)
	const route = "/users/:uid/courses/:cid/group"
	router.Add(http.MethodGet, route, web.GetGroupByUserAndCourse(db))
	// add the route to handler
	requestURL := "/users/" + strconv.FormatUint(user1.ID, 10) + "/courses/" + strconv.FormatUint(course.ID, 10) + "/group"
	r := httptest.NewRequest(http.MethodGet, requestURL, nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)
	router.Find(http.MethodGet, requestURL, c)
	// invoke the prepared handler
	if err := c.Handler()(c); err != nil {
		t.Error(err)
	}
	assertCode(t, w.Code, http.StatusFound)

	var respGroup models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &respGroup); err != nil {
		t.Fatal(err)
	}

	dbGroup, err := db.GetGroup(group.ID)
	if err != nil {
		t.Fatal(err)
	}
	// see models.Group; enrollment field is not transmitted over http
	// we simply ignore enrollments
	dbGroup.Enrollments = nil

	if !reflect.DeepEqual(&respGroup, dbGroup) {
		t.Errorf("have response group %+v, while database has %+v", &respGroup, dbGroup)
	}
}
