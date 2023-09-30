// Grugq-Says
// A simple Go program that uses the OpenAI API to generate explanations for quotes from The Grugq.
// Version: 0.1.0
// Daniel Cuthbert @dcuthbert

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	Green = "\033[32m"
	Blue  = "\033[34m"
	Reset = "\033[0m"
)

func printColored(color, title, content string) {
	fmt.Println(color+title+Reset, content)
	fmt.Println() // this isnt right but it works, dont judge me
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	quotes := map[string]string{
		"The Grugq":                                  "Give a man an 0day and he'll have access for a day, teach a man to phish and he'll have access for life",
		"You can't fight a meme with an exploit":     "The Grugq",
		"Cyber warfare isn't chess, it's calvinball": "The Grugq",
		"An APT is not a toolchain. You can't download your way to parity with Ft Meade":                                                            "The Grugq",
		"OPSEC is a process, not a tool":                                                                                                            "The Grugq",
		"Finding the right level of paranoia is an operational challenge":                                                                           "The Grugq",
		"0days are offensive security by obscurity":                                                                                                 "The Grugq",
		"Just as fragile for attackers as “security by obscurity” is for defenders":                                                                 "The Grugq",
		"If you want to conceal something, don’t swear people to silence, tell as many alternative stories as possible.":                            "The Grugq",
		"I think the American way of cyberwar is: “it is statistically impossible to make mistakes 1% of the time, plus law of large numbers, so…”": "The Grugq",
		"Grugq’s law is: don’t attribute to exploits what can adequately be explained by password theft.":                                           "The Grugq",
		"The P in APT doesn’t stand for “pathetic”":                                                                                                 "The Grugq",
		"Relying on attacker incompetence is no way to go through life":                                                                             "The Grugq",
		"Offensive cyber’s real strategic (i.e., continuing) advantage is a “true positive” success signal. Defenders must deal with this":          "The Grugq",
		"Only break one law at a time.":                                                                                                             "The Grugq",
		"Never lie by accident":                                                                                                                     "The Grugq",
		"ProTip: you’re not worth an 0day":                                                                                                          "The Grugq",
		"Fear of 0day is like being terrified of ninjas instead of cardiovascular disease":                                                          "The Grugq",
		"I’m not going to advise you on how to break the law other than to suggest that you shouldn’t":                                              "The Grugq",
		"Cyber is really only effective as an offensive capability. Defence has mitigation, detection, resilience, etc...but at the end of the day, cyber is a domain that favours the offensive (of course, once on someone else's network, you're on the defensive)": "The Grugq",
		"Make compromises: cost more; yield less; harder to use; easier to find. Analyze them, & stay awake":                                           "The Grugq",
		"Fetishising 0day means that people think once a vulnerability is public there's some sort of automagic immunity":                              "The Grugq",
		"It is surprising how critical good phishing technique is with these APT attacks. Effective phishing is more important than 0day.":             "The Grugq",
		"I think I understand the US strategy against Chinese APT. It’s to flood the APT with so much data they won’t have analysts to review it all.": "The Grugq",
		"The APT that can be named is not the real APT. The way of APT is vast and unknowable. The APT is everywhere & nowhere":                        "The Grugq",
		"APT: repeatable success, interchangeable operators of low to mediocre skill. Easy to train techniques. Consistent results.":                   "The Grugq",
		"Limit the number of people involved to the bare minimum.":                                                                                     "The Grugq",
	}

	randomQuote := getRandomQuote(quotes)
	// The maxTokens value is the maximum number of tokens to generate.
	// The higher the value, the more verbose the response.
	// 100 seems to be the sweet spot
	enrichment := generateExplanation(randomQuote, "davinci", 100)

	// we need to use the printcentre function to make it look pretty
	enrichment = wrapEnrichment(enrichment)

	printColored(Green, "The Grugq Says:", randomQuote)
	printColored(Blue, "We mortals can interpret it as:", enrichment)

}

// wrapEnrichment takes a string as input, and returns the string wrapped in a specific format.
func wrapEnrichment(enrichment string) string {
	// Get the terminal width.
	terminalWidth, _, err := term.GetSize(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Calculate the number of spaces to indent the response.
	padding := terminalWidth/2 - len(enrichment)/2

	// Print the response in a specific format.
	return fmt.Sprintf("%*s%s", padding, "", enrichment)
}

// you need to set the OPENAI_API_KEY environment variable to your API key
// friends dont let friends stick it in the code, kapish? :)

// getRandomQuote takes a map as input, and returns a random quote from that map.

func getRandomQuote(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys[rand.Intn(len(keys))]
}

// you need to set the OPENAI_API_KEY environment variable to your API key
// friends dont let friends stick it in the code, kapish? :)
// I also probably went a bit mad with the error handling, but YOLO

func generateExplanation(quote string, model string, maxTokens int) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		printColored(Reset, "Error: OPENAI_API_KEY environment variable is not set.", "")
		return ""
	}

	request := struct {
		Prompt    string `json:"prompt"`
		MaxTokens int    `json:"max_tokens"`
	}{
		Prompt:    "Can you explain what this really means? " + quote,
		MaxTokens: maxTokens,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		printColored(Reset, "Error marshalling request:", err.Error())
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/engines/"+model+"/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		printColored(Reset, "Error creating HTTP request:", err.Error())
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printColored(Reset, "Error sending HTTP request:", err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		printColored(Reset, "Error: Received HTTP status "+string(rune(resp.StatusCode))+" from OpenAI API.", "")
		return ""
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printColored(Reset, "Error reading response body:", err.Error())
		return ""
	}

	var response ChatGPTResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		printColored(Reset, "Error unmarshalling response:", err.Error())
		return ""
	}

	if len(response.Choices) == 0 || response.Choices[0].Text == "" {
		printColored(Reset, "Error: Received empty answer from OpenAI API.", "")

		return ""
	}

	return response.Choices[0].Text

}
