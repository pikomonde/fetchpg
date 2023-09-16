package fetchpg230715notice

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pikomonde/fetchpg/tools"
)

const (
	solveCapcthaMethodManual = "CAPTCHA_MANUAL"
	solveCapcthaMethodAuto   = "CAPTCHA_AUTO"
)

// config
const prefix = "fetchpg230715notice/output/2023-08/"
const solveCapcthaMethod = solveCapcthaMethodAuto
const apiKey2Captcha = "b7dae04298c886ab818ec47258bcc242"
const startDate = "08/02/2023"
const endDate = "08/19/2023"

func Fetch() {
	fetchFrom(-1, -1, -1)
	// fetchFrom(76, 5, -1)
}

// -1 value for startPage and startRecordInPage for getting all data
func fetchFrom(startCounty, startPage, startRecordInPage int) {
	fmt.Println("---------------------------------")
	client := api2captcha.NewClient(apiKey2Captcha)

	wsURL := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(wsURL).MustConnect()
	defer browser.MustClose()

	// create a new page
	page := browser.MustPage("https://www.georgiapublicnotice.com").MustWaitStable()
	time.Sleep(30 * time.Second)
	/*
		YOU NEED TO INSTALL OPEN VPN IN THIS 30 SECOND WINDOW, MANUALLY:
		- OPEN GOOGLE, SEARCH FOR "veepn", INSTALL EXTENSION;
		- TURN ON THE VPN;
		- CLOSE OTHER TABS, BACK TO PUBLIC NOTICE TAB;
		- WAIT.
	*/
	tools.Screenshot(page, page.MustElement("body"), prefix+"01_start.png")
	tools.Screenshot(page, page.MustElement("body"), prefix+"01_start.png")

	// search input
	page.MustElement(".main-search .keyword-search .wsForm .form_item [id$='txtSearch']").MustInput("security deed")

	// county filter (getting list)
	counties := page.MustElements(".main-search .refine-publication .container [id$='divCounty'] #countyDiv li")

	// date input
	page.MustElement(".main-search .refine-publication .container [id$='divDateRange'] .header a").MustClick()
	page.MustWaitStable()
	page.MustWaitStable()
	page.MustElement(".main-search .refine-publication .container [id$='divDateRange'] [id$='txtDateFrom']").MustSelectAllText().MustInput(startDate)
	page.MustElement(".main-search .refine-publication .container [id$='divDateRange'] [id$='txtDateTo']").MustSelectAllText().MustInput(endDate)
	page.MustElement(".main-search .refine-publication .container [id$='divDateRange'] [id$='rbRange']").MustClick()
	page.MustWaitStable()
	page.MustWaitStable()
	tools.Screenshot(page, page.MustElement("body"), prefix+"02_search_input.png")
	tools.Screenshot(page, page.MustElement("body"), prefix+"02_search_input.png")

	// search
	page.MustElement(".main-search .buttons [id$='btnGo']").MustClick()
	waitPageLoad(page, true, 1*time.Second, "[id$='lnkResults']")

	// loop counties
	prevIdxCounty := -1
	for idxCounty := range counties {
		// ===== hardcode here to continue =====
		if startCounty >= 0 {
			if idxCounty+1 < startCounty {
				continue
			}
		}

		// county filter
		page.MustElement(".main-search .refine-publication .container [id$='divCounty'] .header a").MustClick()
		page.MustWaitStable()
		page.MustWaitStable()
		counties2 := page.MustElements(".main-search .refine-publication .container [id$='divCounty'] #countyDiv li")
		countyName := strings.ReplaceAll(strings.ToLower(counties2[idxCounty].MustElement("label").MustText()), " ", "_")
		fmt.Printf("start fetching for county %s (%d/%d)\n", countyName, idxCounty+1, len(counties))
		if prevIdxCounty >= 0 {
			// deselect previous
			counties2[prevIdxCounty].MustClick()
			page.MustWaitStable()
			page.MustWaitStable()
		}
		// select current
		counties2 = page.MustElements(".main-search .refine-publication .container [id$='divCounty'] #countyDiv li")
		counties2[idxCounty].MustClick()
		page.MustWaitStable()
		page.MustWaitStable()

		// search
		page.MustElement(".main-search .buttons [id$='btnGo']").MustClick()
		waitPageLoad(page, true, 1*time.Second, "[id$='lnkResults']")
		tools.Screenshot(page, page.MustElement("body"), prefix+fmt.Sprintf("03_search_page_%s.png", countyName))
		tools.Screenshot(page, page.MustElement("body"), prefix+fmt.Sprintf("03_search_page_%s.png", countyName))

		// init county data
		dataCounty := NewData()

		// getting page info
		totalPage := 1
		if len(page.MustElements(".pager")) > 0 {
			// curPage = page.MustElement(".pager [id$='lblCurrentPage']").MustText()
			totalPageText := strings.TrimSpace(page.MustElement(".pager [id$='lblTotalPages']").MustText())
			// strings.Contains(totalPageText,curPage) // alternative to compare page total
			totalPageArr := strings.Split(totalPageText, " ")
			if len(totalPageArr) == 3 {
				totalPageStr := totalPageArr[1]
				totalPage, _ = strconv.Atoi(totalPageStr)
			}
		}

		// loop pages
		for curPage := 1; curPage <= totalPage; curPage++ {
			// ===== hardcode here to continue =====
			if startPage >= 0 {
				if curPage < startPage {
					nextPage(page)
					continue
				}
			}

			// loop records
			records := page.MustElements("input[id$='btnView2']")
			for idxRecord := range records {
				// ===== hardcode here to continue =====
				if startRecordInPage >= 0 {
					if idxRecord+1 < startRecordInPage {
						continue
					}
				}

				records2 := page.MustElements("input[id$='btnView2']")
				records2[idxRecord].MustClick()
				waitPageLoad(page, false, 1*time.Second, "")

				// solve captcha
				if len(page.MustElements("[id$='pnlReCaptcha']")) > 0 {
					// notify user if manual
					if solveCapcthaMethod == solveCapcthaMethodManual {
						// exec.Command("say captcha captcha captcha captcha captcha captcha captcha captcha captcha captcha captcha").
						// 	Run()
					}

					solveCaptcha(page, client)
				}

				// fetching data
				content := page.MustElement(".notice [id$='pnlNoticeContent'] [id$='lblContentText']").MustText()
				// pubName := strings.ToLower(strings.ReplaceAll(strings.Split(page.MustElement("#detail [id$='lblPubName']").MustText(), "-")[0], " ", "_"))
				postID := ""
				urlSplitPostID := strings.Split(page.MustInfo().URL, "&ID=")
				if len(urlSplitPostID) == 2 {
					postID = urlSplitPostID[1]
				}

				fmt.Printf("start fetching for county %s (%d/%d), page (%d/%d), record in page (%d/%d), postID: %s\n",
					countyName, idxCounty+1, len(counties), curPage, totalPage, idxRecord+1, len(records), postID)

				// screenshot
				fname := fmt.Sprintf("screenshot/%s-%s.png", countyName, postID)
				tools.Screenshot(page, page.MustElement("body"), prefix+fname)
				tools.Screenshot(page, page.MustElement("body"), prefix+fname)

				// saving county & final data
				dataCounty.Add(page.MustInfo().URL, countyName, fname, content)
				dataCounty.SaveFile(prefix + fmt.Sprintf("output_%s.csv", countyName))

				// back
				page.MustElement(".backlink a").MustClick()
				waitPageLoad(page, true, 1*time.Second, "[id$='lnkResults']")

				// if reset to homepage
				if len(page.MustElements("#capitol")) > 0 {
					gotoCurPageFromHomePage(page, curPage)
				}
			}
			startRecordInPage = -1

			// next page
			if curPage < totalPage {
				nextPage(page)
			}
		}
		prevIdxCounty = idxCounty
		startPage = -1

		// saving county csv
		dataCounty.SaveFile(prefix + fmt.Sprintf("output_%s.csv", countyName))
	}

}

func nextPage(page *rod.Page) {
	page.MustElement(".pager [id$='btnNext']").MustClick()
	waitPageLoad(page, true, 1*time.Second, "[id$='lnkResults']")
}

func gotoCurPageFromHomePage(page *rod.Page, curPage int) {
	fmt.Println("page reset, setup to current page")

	page.MustElement(".main-search .buttons [id$='btnGo']").MustClick()
	waitPageLoad(page, true, 1*time.Second, "[id$='lnkResults']")

	// to current page
	for idxSetupCurPage := 1; idxSetupCurPage < curPage; idxSetupCurPage++ {
		nextPage(page)
	}

	fmt.Println("success setup to current page")
}

func solveCaptcha(page *rod.Page, cli *api2captcha.Client) {

	fmt.Println("solving captcha")

	if solveCapcthaMethod == solveCapcthaMethodManual {

		tools.Notify()
		recapcthaFrame := page.MustElement("[id$='pnlReCaptcha']").MustFrame()
		page.MustSetWindow(0, 0, *recapcthaFrame.MustGetWindow().Width, *recapcthaFrame.MustGetWindow().Height)
		el := page.MustElement("[id$='pnlReCaptcha']")
		el.MustScrollIntoView()

		// next page
		// el.MustWaitInvisible()
		_ = el.WaitInvisible() // there will be error because it change page and the element is destroyed

	} else if solveCapcthaMethod == solveCapcthaMethodAuto {

		// get siteKey
		siteKeyEl := page.MustElement("#recaptcha")
		siteKey := *siteKeyEl.MustAttribute("data-sitekey")
		fmt.Println("recaptcha siteKey:", siteKey)

		// hit 2captcha
		cap := api2captcha.ReCaptcha{
			SiteKey:   siteKey,
			Url:       page.MustInfo().URL,
			Invisible: true,
			Action:    "verify",
		}
		req := cap.ToRequest()
		now := time.Now()
		fmt.Println("2captcha process:", now)
		code, err := cli.Solve(req)
		if err != nil {
			fmt.Println("ERROR Recaptcha:", err)
			tools.Notify()
		}
		fmt.Println("2captcha response code:", code)
		balance, _ := cli.GetBalance()
		fmt.Println("2captcha finished:", time.Now().Sub(now).Seconds())
		fmt.Println("current balance:", balance)

		// put 2captcha response to recaptcha
		page.MustWaitElementsMoreThan("#g-recaptcha-response", 0)
		page.MustEval(`() => (document.querySelector("#g-recaptcha-response").style = "")`)
		page.MustElement("#g-recaptcha-response").MustSelectAllText().MustInput(code)

		// next page
		page.MustElement("[id$='btnViewNotice']").MustClick()
		waitPageLoad(page, false, 1*time.Second, "")

	}

	fmt.Println("success solving captcha")
}

func waitPageLoad(page *rod.Page, isWaitPageLoad bool, waitDuration time.Duration, selector string) {
	time.Sleep(waitDuration)
	if selector != "" {
		page.MustWaitElementsMoreThan(selector, 0)
	}
	// if isWaitPageLoad {
	// 	page.MustWaitLoad()
	// }
	page.MustWaitStable()
	page.MustWaitStable()
}
