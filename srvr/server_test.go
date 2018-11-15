package srvr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TestGetAll Fetches all entries
func TestGetAll(t *testing.T) {
	fmt.Println("Testing GetAll")
	// Create a request to serveHTTP. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1 := httptest.NewRequest("GET", "/books", nil)
	req2 := httptest.NewRequest("GET", "/book", nil)
	rqsts := []*http.Request{req1, req2}
	statusOutArr := []int{http.StatusOK, http.StatusNotFound}
	getResponse(t, rqsts, statusOutArr)

}

//GetTest : Fetches a specified entry
func TestGet(t *testing.T) {
	fmt.Println("")
	fmt.Println("Testing GET")
	// Create a request to serveHTTP. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1 := httptest.NewRequest("GET", "/books/1", nil)
	req2 := httptest.NewRequest("GET", "/books/30", nil)
	req3 := httptest.NewRequest("GET", "/books/sfsg", nil)
	req4 := httptest.NewRequest("GET", "/books/", nil)
	rqsts := []*http.Request{req1, req2, req3, req4}
	statusOutArr := []int{http.StatusOK, http.StatusNoContent, http.StatusBadRequest, http.StatusNotFound}
	getResponse(t, rqsts, statusOutArr)
}

//TestPost
func TestPost(t *testing.T) {
	fmt.Println("")
	fmt.Println("Testing POST")
	//First create a few test POST body
	var bks [4]Book
	bks[0] = Book{ID: "3", Name: "Catch 22", Author: "Joseph Heller", Count: 43}
	bks[1] = Book{ID: "3", Name: "Catch 22", Author: "Joseph Heller", Count: 43}
	bks[2] = Book{ID: "", Name: "Get Out", Author: "John Doe", Count: 4}
	bks[3] = Book{ID: "5", Name: "Hamlet", Author: "William Shakespeare", Count: 3}
	//create a slice of 4 requests
	var rqsts []*http.Request
	rqsts = make([]*http.Request, 4)
	// 4 url to api endpoint, last one is incorrect
	reqURL := []string{"/books", "/books", "/books", "/book"}
	//populate request array
	for i := 0; i < 4; i++ {
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(bks[i])
		// Create a request to pass to our handler. We have a POST to make, so we'll
		// pass body as the third parameter.
		rqsts[i] = httptest.NewRequest("POST", reqURL[i], body)
	}

	statusOutArr := []int{http.StatusOK, http.StatusConflict, http.StatusBadRequest, http.StatusNotFound}
	getResponse(t, rqsts, statusOutArr)

}

//TestUpdate
func TestUpdate(t *testing.T) {
	fmt.Println("")
	fmt.Println("Testing Update")
	//First create a few test POST body
	var bks [4]Book
	bks[0] = Book{ID: "5", Name: "Hamlet", Author: "William Shakespeare", Count: 3}
	bks[1] = Book{ID: "6", Name: "Catch 22", Author: "Joseph Heller", Count: 43}
	bks[2] = Book{ID: "", Name: "Get Out", Author: "John Doe", Count: 4}
	bks[3] = Book{ID: "5", Name: "Hamlet", Author: "William Shakespeare", Count: 3}
	//create a slice of 4 requests
	var rqsts []*http.Request
	rqsts = make([]*http.Request, 4)
	// 4 url to api endpoint, last three is incorrect
	reqURL := []string{"/books/1", "/books/7", "/books/hj", "/book"}
	//populate request array
	for i := 0; i < 4; i++ {

		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(bks[i])
		// Create a request to pass to our handler. We have a POST to make, so we'll
		// pass body as the third parameter.
		rqsts[i] = httptest.NewRequest("PUT", reqURL[i], body)

	}

	// Create a request to pass to our handler. We have a POST to make, so we'll
	// pass body as the third parameter.

	statusOutArr := []int{http.StatusOK, http.StatusNotFound, http.StatusBadRequest, http.StatusNotFound}
	getResponse(t, rqsts, statusOutArr)

}

//DeleteTest : Deletes a specified entry
func TestDelete(t *testing.T) {
	fmt.Println("")
	fmt.Println("Testing DELETE")
	// Create a request to serveHTTP. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1 := httptest.NewRequest("DELETE", "/books/2", nil)
	req2 := httptest.NewRequest("DELETE", "/books/2", nil)
	req3 := httptest.NewRequest("DELETE", "/books/sfsg86", nil)
	req4 := httptest.NewRequest("DELETE", "/books/", nil)
	rqsts := []*http.Request{req1, req2, req3, req4}
	statusOutArr := []int{http.StatusOK, http.StatusGone, http.StatusBadRequest, http.StatusNotFound}
	getResponse(t, rqsts, statusOutArr)
}

//getResponse takes https requests, and expected output status
//outputs error message or print successfully returned messeges
func getResponse(t *testing.T, req []*http.Request, statusOutArr []int) {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.

	for i, val := range req {
		fmt.Println("Test case : ", i+1)
		res := httptest.NewRecorder()

		val.Header.Set("Authorization", "Base "+base64.StdEncoding.EncodeToString([]byte("tom95:pass1")))
		//Now the response request pair is served via http
		router.ServeHTTP(res, val)

		// Check the status code is what we expect.
		if statusOutArr[i] != res.Code {
			fmt.Println("Fail.")
			t.Error("Handler returned wrong status code: got ", res.Code, "want ", statusOutArr[i])

		} else {
			fmt.Println("Pass.  {{", res.Code, "}}")
			if res.Code == http.StatusOK {
				printStruct(res.Body.String())
				fmt.Println("")
			}
		}

	}
	fmt.Println("")

}

func printStruct(responseBody string) {
	var books []Book
	json.Unmarshal([]byte(responseBody), &books)
	fmt.Println("Response body returned: ")

	for _, bk := range books {
		fmt.Println("	Book ID = ", bk.ID)
		fmt.Println("	Name = ", bk.Name)
		fmt.Println("	Author = ", bk.Author)
		fmt.Println("	Remaining", bk.Count)
		fmt.Println(" ")
	}

}

/*


 */
