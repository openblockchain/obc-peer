/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package crypto

import (
	"errors"
	"time"
	"sort"
	"encoding/hex"
	utils "github.com/hyperledger/fabric/core/crypto/utils"
)

// The Multi-threaded tCertPool is currently not used.
// It plays only a role in testing.
type tCertPoolMultithreadingImpl struct {
	client *clientImpl

	tCertChannel         chan tCert
	tCertChannelFeedback chan struct{}
	done                 chan struct{}
}

func (tCertPool *tCertPoolMultithreadingImpl) Start() (err error) {
	// Start the filler
	go tCertPool.filler()

	return
}

func (tCertPool *tCertPoolMultithreadingImpl) Stop() (err error) {
	// Stop the filler
	tCertPool.done <- struct{}{}

	// Store unused TCert
	tCertPool.client.debug("Store unused TCerts...")

	tCerts := []tCert{}
	for {
		if len(tCertPool.tCertChannel) > 0 {
			tCerts = append(tCerts, <-tCertPool.tCertChannel)
		} else {
			break
		}
	}

	tCertPool.client.debug("Found %d unused TCerts...", len(tCerts))

//TODO FIX MT POOL
	//tCertPool.client.ks.storeUnusedTCerts(tCerts)

	tCertPool.client.debug("Store unused TCerts...done!")

	return
}



func (tCertPool *tCertPoolMultithreadingImpl) CalculateAttributesHash(attributes map[string]string) (attrHash string) {
	
	keys := make([]string, len(attributes))
	 
	 for k := range attributes {
        keys = append(keys, k)
    }
	 
	  sort.Strings(keys)
	 
	  
	  values := make([]byte, len(keys))
	  
	 for _,k := range keys {
	 	
	 	vb := []byte(k)
	 	for _,bval := range vb {
	 		 values = append(values, bval)
	 	}
	 	
	 	vb = []byte(attributes[k])
	 	for _,bval := range vb {
	 		 values = append(values, bval)
	 	}
    }
	
	
	attributesHash := utils.Hash(values)
	
	return  hex.EncodeToString(attributesHash)
	
}


func (tCertPool *tCertPoolMultithreadingImpl) GetNextTCert(attributes map[string]string) (tCertBlock *TCertBlock, err error) {
	for i := 0; i < 3; i++ {
		tCertPool.client.debug("Getting next TCert... %d out of 3", i)
		
		tCertBlock = new(TCertBlock)
		
		select {
		case tCertBlock.tCert = <-tCertPool.tCertChannel:
			break
		case <-time.After(30 * time.Second):
			tCertPool.client.error("Failed getting a new TCert. Buffer is empty!")

			//return nil, errors.New("Failed getting a new TCert. Buffer is empty!")
		}
		if tCertBlock.tCert != nil {
			// Send feedback to the filler
			tCertPool.tCertChannelFeedback <- struct{}{}
			break
		}
	}

	if tCertBlock.tCert == nil {
		// TODO: change error here
		return nil, errors.New("Failed getting a new TCert. Buffer is empty!")
	}

	tCertPool.client.debug("Cert [% x].", tCertBlock.tCert.GetCertificate().Raw)

	// Store the TCert permanently
	tCertPool.client.ks.storeUsedTCert(tCertBlock)

	tCertPool.client.debug("Getting next TCert...done!")

	return
}

func (tCertPool *tCertPoolMultithreadingImpl) AddTCert(tCertBlock *TCertBlock) (err error) {
	tCertPool.client.debug("New TCert added.")
	tCertPool.tCertChannel <- tCertBlock.tCert

	return
}

func (tCertPool *tCertPoolMultithreadingImpl) init(client *clientImpl) (err error) {
	tCertPool.client = client

	tCertPool.tCertChannel = make(chan tCert, client.conf.getTCertBatchSize()*2)
	tCertPool.tCertChannelFeedback = make(chan struct{}, client.conf.getTCertBatchSize()*2)
	tCertPool.done = make(chan struct{})

	return
}

func (tCertPool *tCertPoolMultithreadingImpl) filler() {
	// Load unused TCerts
	stop := false
//	full := false //TODO FIX
	for {
		// Check if Stop was called
		select {
		case <-tCertPool.done:
			tCertPool.client.debug("Force stop!")
			stop = true
		default:
		}
		if stop {
			break
		}

		tCertDER, err := tCertPool.client.ks.loadUnusedTCert()
		if err != nil {
			tCertPool.client.error("Failed loading TCert: [%s]", err)
			break
		}
		if tCertDER == nil {
			tCertPool.client.debug("No more TCerts in cache!")
			break
		}

//TODO FIX
		/*tCert, err := tCertPool.client.getTCertFromDER(tCertDER)
		if err != nil {
			tCertPool.client.error("Failed paring TCert [% x]: [%s]", tCertDER, err)

			continue
		}


		// Try to send the tCert to the channel if not full
		select {
		case tCertPool.tCertChannel <- tCert:
			tCertPool.client.debug("TCert send to the channel!")
		default:
			tCertPool.client.debug("Channell Full!")
			full = true
		}
		if full {
			break
		}*/
	}

	tCertPool.client.debug("Load unused TCerts...done!")

	if !stop {//TODO FIXcl
		/*ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-tCertPool.done:
				stop = true
				tCertPool.client.debug("Done signal.")
			case <-tCertPool.tCertChannelFeedback:
				tCertPool.client.debug("Feedback received. Time to check for tcerts")
			case <-ticker.C:
				tCertPool.client.debug("Time elapsed. Time to check for tcerts")
			}

			if stop {
				tCertPool.client.debug("Quitting filler...")
				break
			}

			if len(tCertPool.tCertChannel) < tCertPool.client.conf.getTCertBatchSize() {
				tCertPool.client.debug("Refill TCert Pool. Current size [%d].",
					len(tCertPool.tCertChannel),
				)

				var numTCerts = cap(tCertPool.tCertChannel) - len(tCertPool.tCertChannel)
				if len(tCertPool.tCertChannel) == 0 {
					numTCerts = cap(tCertPool.tCertChannel) / 10
				}

				tCertPool.client.info("Refilling [%d] TCerts.", numTCerts)

				err := tCertPool.client.getTCertsFromTCA(nil, nil, numTCerts)
				if err != nil {
					tCertPool.client.error("Failed getting TCerts from the TCA: [%s]", err)
				}
			}
		}*/
	}

	tCertPool.client.debug("TCert filler stopped.")
}
