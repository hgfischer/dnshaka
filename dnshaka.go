// Copyright 2017 Herbert Fischer. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// DNSHaka is an DNS benchmarking tool.
package main

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/tylertreat/bench"
)

func main() {
	r := &dnsRequesterFactory{
		Address:          "8.8.8.8:53",
		Question:         dns.Question{Name: "hgfdishas.asd.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
		RecursionDesired: true,
	}

	benchmark := bench.NewBenchmark(r, 1000, 1, 10*time.Second, 0)
	summary, err := benchmark.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(summary)
	summary.GenerateLatencyDistribution(nil, "dns.txt")
}

type dnsRequesterFactory struct {
	Address          string
	Question         dns.Question
	RecursionDesired bool
}

func (f *dnsRequesterFactory) GetRequester(num uint64) bench.Requester {
	return &dnsRequester{
		Address:          f.Address,
		Question:         f.Question,
		RecursionDesired: f.RecursionDesired,
	}
}

type dnsRequester struct {
	Address          string
	Question         dns.Question
	RecursionDesired bool
	msg              *dns.Msg
}

func (r *dnsRequester) Setup() error {
	r.msg = new(dns.Msg)
	r.msg.Id = dns.Id()
	r.msg.RecursionDesired = r.RecursionDesired
	r.msg.Question = make([]dns.Question, 1)
	r.msg.Question[0] = r.Question
	return nil
}

func (r *dnsRequester) Request() error {
	c := new(dns.Client)
	_, _, err := c.Exchange(r.msg, r.Address)
	return err
}

func (r *dnsRequester) Teardown() error {
	return nil
}
