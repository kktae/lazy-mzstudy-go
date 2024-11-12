package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var (
	MegazoneId string
	MegazonePw string

	MzStudyUrl = url.URL{Scheme: "https", Host: "mz.livestudy.com"}

	// mutex 선언
	globalMutex sync.Mutex
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	MegazoneId = viper.GetString("mz_id")
	MegazonePw = viper.GetString("mz_pw")
	if MegazoneId == "" {
		log.Fatal("https://mz.livestudy.com 에서 사용하는 ID를 정확하게 입력해주세요.")
	}
	if MegazonePw == "" {
		log.Fatal("https://mz.livestudy.com 에서 사용하는 Password를 정확하게 입력해주세요.")
	}
}

func SetClientMzStudyBaseUrl(client *resty.Client) {
	client.SetBaseURL(MzStudyUrl.String())
}

func SetClientHeaders(client *resty.Client) {
	client.SetHeaders(map[string]string{
		"Accept": "application/json",
		// "Accept-Encoding": "gzip, deflate, br",
		"Connection": "keep-alive",
		// "Content-Type":    "application/x-www-form-urlencoded",
		"Origin":     "https://mz.livestudy.com",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	})
}

func SetClientRetries(client *resty.Client) {
	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})
}

func SetClientMzStudySessionId(client *resty.Client, jsessionId string) {
	if jsessionId == "" {
		jsessionId = GetJSESSIONId()
	}
	if !strings.Contains(jsessionId, ".mz_was1") {
		jsessionId += ".mz_was1"
	}

	client.SetCookie(&http.Cookie{
		Name:     "JSESSIONID",
		Value:    jsessionId,
		Domain:   "mz.livestudy.com",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
}

func LoginMzStudy(client *resty.Client) (*resty.Response, error) {
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
			"Content-Type":    "application/x-www-form-urlencoded",
			"Origin":          "https://mz.livestudy.com",
			"Referer":         "https://mz.livestudy.com/login",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		}).
		SetFormData(map[string]string{
			"info_os":           "Windows 10",
			"info_browser":      "Chrome_113",
			"info_resolution":   "1920x1056",
			"info_flashVersion": "0,0,0",
			"j_username":        MegazoneId,
			"password":          MegazonePw,
			"j_password":        EncPW(MegazonePw),
		}).
		Post("/j_spring_security_check")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 나의 강의실에서 학습현황 받아오기
func RequestGetMyPage(client *resty.Client) (*resty.Response, error) {
	resp, err := client.R().Get("/mypage/main")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 학습현황 파싱
func ParseCourseIdsFromMyPage(body string) []string {
	dataCourses := []string{}

	// HTML 문자열에서 goquery 문서 생성
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	// name 속성이 "btnLesson"이고 href 속성이 "#"인 모든 요소를 선택하고 data-courseid 값을 가져옴
	doc.Find("[href='#'][name='btnLesson']").Each(func(i int, s *goquery.Selection) {
		courseID, exists := s.Attr("data-courseid")
		if exists {
			// fmt.Println("data-courseid:", courseID)
			dataCourses = append(dataCourses, courseID)
		}
	})
	return dataCourses
}

func RequestGetCourseData(client *resty.Client, courseId string) (*Response, error) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	var err error
	var result Response

	_, err = client.R().
		SetQueryString(fmt.Sprintf("courseId=%s", courseId)).
		SetResult(&result).
		Get(fmt.Sprintf("/proxy/attendances/%s.json", courseId))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func RequestGetAttendances(client *resty.Client, courseId int, lessonId int) (*Response, error) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	var err error
	var result Response

	_, err = client.R().
		SetQueryString(fmt.Sprintf("courseId=%d&courseLessonId=%d", courseId, lessonId)).
		SetResult(&result).
		Get("/proxy/attendances.json")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func RequestPostUpdateAttend(client *resty.Client, courseId int, lessonId int, finishedPage int, isReset bool, logId int) (*Response, error) {
	// globalMutex.Lock()
	// defer globalMutex.Unlock()

	var err error
	var result Response

	_, err = client.R().
		SetFormData(map[string]string{
			"finishedPage":          strconv.Itoa(finishedPage),
			"resolvedQuestionCount": "0",
			"isReset":               strconv.FormatBool(isReset), // origin: true
			// "user.id":               strconv.Itoa(0),             // no need this field
			"id": strconv.Itoa(logId),
		}).
		SetResult(&result).
		Post(fmt.Sprintf("/proxy/course/%d/lessons/%d/updateAttend.json", courseId, lessonId))
	if err != nil {
		return nil, err
	}
	if result.Status != "SUCCEED" {
		return nil, fmt.Errorf(result.Message.(string))
	}
	return &result, nil
}

func RequestStudyLesson(client *resty.Client, courseId, courceLessonId, totalPageCount int, requiredTotalLearningTime, sleepDuration time.Duration) {

	wg := sync.WaitGroup{}

	// courseId := 1015
	// courceLessonId := 21311
	// totalPageCount := 5

	// // 총 학습 시간 (분)
	// totalLearningTime := time.Minute * 20

	// // 대기 시간 (초)
	// sleepDuration := time.Second * 60

	// 생성할 스레드 개수 // (초)를 기준으로 계산
	threadCount := int(math.Ceil(requiredTotalLearningTime.Seconds() / sleepDuration.Seconds()))

	for i := 0; i < threadCount; i++ {

		modPage := (i % totalPageCount) + 1

		if totalPageCount < (i + 1) {
			modPage = 1
		}

		wg.Add(1) // 각 Goroutine 시작 시 WaitGroup 증가
		// 페이지별로 goroutine 실행
		go func(page int) {
			defer wg.Done() // 각 Goroutine 종료 시 WaitGroup 감소

			// 현재 시간을 기준으로 begin
			resp, err := RequestPostUpdateAttend(client, courseId, courceLessonId, page, false, 0)
			if err != nil {
				log.Println("begin update attend error:", err)
				return
			}
			if resp.HasErrors {
				log.Println("begin update attend response has error:", resp.HasErrors)
				return
			}

			// 특정 시간 대기
			time.Sleep(sleepDuration + (time.Second * 5))

			// update request 성공 시, 대기한 시간만큼 기록됨
			resp, err = RequestPostUpdateAttend(client, courseId, courceLessonId, page, false, resp.Data.LogID)
			if err != nil {
				log.Println("end update attend error:", err)
				return
			}
			if resp.HasErrors {
				log.Println("end update attend response has error:", resp.HasErrors)
				return
			}

			resp, err = RequestPostUpdateAttend(client, courseId, courceLessonId, page, false, resp.Data.LogID)
			if err != nil {
				log.Println("after update attend error:", err)
				return
			}
			if resp.HasErrors {
				log.Println("after update attend response has error:", resp.HasErrors)
				return
			}

		}(modPage)

		// Rate limit 방지를 위한 잠시 대기
		src := rand.NewPCG(rand.Uint64(), rand.Uint64())
		rng := rand.New(src)
		randomMillis := rng.IntN(501) + 1100
		time.Sleep(time.Millisecond * time.Duration(randomMillis))
	}

	// 모든 Goroutine이 종료될 때까지 대기
	wg.Wait()
}

func main() {

	log.Println("start process")
	defer func() {
		log.Println("end process")
	}()

	loadConfig()

	client := resty.New()

	// client.SetDebug(true)
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

	SetClientMzStudyBaseUrl(client)
	SetClientHeaders(client)
	SetClientRetries(client)
	SetClientMzStudySessionId(client, "")

	var resp *resty.Response
	var err error

	_ = resp
	_ = err

	_, err = LoginMzStudy(client)
	if err != nil {
		log.Println(err)
		return
	}

	// /* logging when debug flag */
	// log.Println("----------------------------------------------")
	// log.Println("------------------------------------------POST")
	// log.Printf("\nError: %v", err)
	// log.Printf("\nResponse Status Code: %v", resp.StatusCode())
	// log.Printf("\nResponse Status: %v", resp.Status())
	// // log.Printf("\nResponse Body: %v", resp)
	// log.Printf("\nResponse Time: %v", resp.Time())
	// // log.Printf("\nCookies: %v", resp.Cookies())
	// log.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	// log.Println("----------------------------------------------")

	// RequestDelayLesson(client)

	resp, err = RequestGetMyPage(client)
	if err != nil {
		log.Println(err)
		return
	}

	courceIds := ParseCourseIdsFromMyPage(string(resp.Body()))
	if len(courceIds) <= 0 {
		log.Println("course id 슬라이스의 길이가 잘못되었습니다.")
		return
	}

	wgCourse := sync.WaitGroup{}
	for _, courceId := range courceIds {
		courceResp, err := RequestGetCourseData(client, courceId)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(courceResp.Data.Course.Attendances) <= 0 {
			log.Println("attendances의 슬라이스 길이가 잘못되었습니다.")
			continue
		}

		wgCourse.Add(1)
		go func() {
			defer wgCourse.Done() // 각 Course Goroutine 종료 시 WaitGroup 감소

			wgLession := sync.WaitGroup{}
			for _, attendance := range courceResp.Data.Course.Attendances {

				if attendance.CourseLesson.ID <= 0 {
					continue
				}

				finishedLearningTimeDur := time.Second * time.Duration(attendance.FinishedLearningTime)
				lessionLearningTimeDur := time.Minute * time.Duration(attendance.CourseLesson.Lesson.LearningTime)
				// log.Println("finishedLearningTimeDur.Seconds():", finishedLearningTimeDur.Seconds())
				// log.Println("lessionLearningTimeDur.Seconds():", lessionLearningTimeDur.Seconds())
				if finishedLearningTimeDur.Seconds() >= lessionLearningTimeDur.Seconds() {
					log.Printf("skip already study lesson [%05d][%s]\n", attendance.CourseLesson.ID, attendance.CourseLesson.Lesson.Title)
					continue
				}

				log.Printf("begin study lesson [%05d][%s]\n", attendance.CourseLesson.ID, attendance.CourseLesson.Lesson.Title)

				attendancesResp, err := RequestGetAttendances(client, courceResp.Data.Course.ID, attendance.CourseLesson.ID)
				if err != nil {
					log.Println(err)
					continue
				}

				// totalLearningTime := time.Minute * time.Duration(attendance.CourseLesson.Lesson.LearningTime)

				leftTotalLearningTime := lessionLearningTimeDur - finishedLearningTimeDur
				sleepDuration := time.Second * 60

				wgLession.Add(1)
				go func() {
					defer wgLession.Done() // 각 Goroutine 종료 시 WaitGroup 감소
					RequestStudyLesson(client, courceResp.Data.Course.ID, attendance.CourseLesson.ID, attendancesResp.Data.Standard.Lesson.TotalPages, leftTotalLearningTime, sleepDuration)
					log.Printf("end study lesson [%05d][%s]\n", attendance.CourseLesson.ID, attendance.CourseLesson.Lesson.Title)
				}()
			}
			wgLession.Wait()
		}()
	}

	wgCourse.Wait()
}
