package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

type RobotChecker struct {
	robot map[string](*robotstxt.RobotsData)
}

func (r *RobotChecker) init() {
	r.robot = map[string](*robotstxt.RobotsData){}
}

func (r *RobotChecker) checkIfAllowed(_url string) bool {

	u, err := url.Parse(_url)

	if err != nil {
		fmt.Println("Error parsing url:", err)
		return false
	}

	if u.Scheme != "https" {
		return false
	}

	robotsUri := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)
	rest := u.Path
	if u.RawQuery != "" {
		rest = rest + "?" + u.RawQuery
	}
	if u.Fragment != "" {
		rest = rest + "#" + u.Fragment
	}

	robot, exists := r.robot[robotsUri]
	if !exists {
		res, err := http.Get(robotsUri)
		if err != nil {
			fmt.Println("Error querying robots.txt:", err)
			return false
		}
		robot, err = robotstxt.FromResponse(res)
		r.robot[robotsUri] = robot
		if err != nil {
			fmt.Println("Error creating robots.txt checker:", err)
			return false
		}
		res.Body.Close()

	}

	group := robot.FindGroup("SearchCrawler Bot")
	allowed := group.Test(rest)

	return allowed
}
