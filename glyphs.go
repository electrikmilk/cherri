/*
 * Copyright (c) Cherri
 */

/*
 * Glyphs sourced from https://github.com/pfgithub/scpl under MIT License.
 * 2019 pfgithub.

 * Under courtesy of https://github.com/OpenJelly/Open-Jellycore.
 * Created by Taylor Lineman on 6/2/23.
 */

package main

var iconGlyph int64 = 61440

var glyphs = map[string]int{
	"car":                        59452,
	"carMultiple":                61446,
	"electricCar":                61447,
	"bus":                        59678,
	"busDouble":                  61448,
	"tram":                       61449,
	"tramTunnel":                 61450,
	"bike":                       59668,
	"motorcycle":                 59783,
	"ambulance":                  59652,
	"airplane":                   59648,
	"sailboat":                   59823,
	"house":                      59755,
	"church":                     59688,
	"buildings":                  59677,
	"shoppingCart":               59828,
	"handbag":                    59750,
	"groceryStore":               59747,
	"utensils":                   59863,
	"fuelstation":                59741,
	"thermometer":                59854,
	"sun":                        59845,
	"moon":                       59782,
	"moonCircle":                 61517,
	"snowflake":                  59835,
	"cloud":                      59714,
	"raincloud":                  59715,
	"umbrella":                   59861,
	"pineTree":                   59731,
	"flower":                     59468,
	"fire":                       59734,
	"footprints":                 59738,
	"signs":                      59724,
	"binoculars":                 59669,
	"compass":                    59717,
	"globe":                      59412,
	"mountain":                   59785,
	"picture":                    59784,
	"filmstrip":                  59733,
	"camera":                     59682,
	"movieCamera":                59402,
	"microphone":                 59780,
	"videoIcon":                  59864,
	"clipboard":                  59711,
	"calendar":                   59681,
	"chatBubble":                 59414,
	"messageBubbles":             59403,
	"textBubble":                 59779,
	"envelope":                   59773,
	"envelopeOpen":               59774,
	"paperAirplane":              59836,
	"paperAirplaneCircle":        61462,
	"briefcase":                  59676,
	"folder":                     59737,
	"creditCard":                 59719,
	"watch":                      59865,
	"phone":                      59814,
	"laptop":                     59436,
	"keyboard":                   59446,
	"keyboardOld":                59494,
	"calculator":                 59680,
	"barGraph":                   59662,
	"printer":                    59817,
	"hardDrive":                  59752,
	"server":                     59722,
	"database":                   59519,
	"networkStorage":             59826,
	"archive":                    59653,
	"cube":                       59721,
	"tv":                         59851,
	"gameController":             59742,
	"puzzlePiece":                59818,
	"headphones":                 59753,
	"headphonesCircle":           61479,
	"ear":                        61481,
	"musicNote":                  59790,
	"volumeLow":                  59839,
	"volumeMid":                  61470,
	"volumeHigh":                 61471,
	"mute":                       61472,
	"speaker":                    61473,
	"hifiSpeaker":                61478,
	"desktopSpeaker":             61474,
	"bookshelf":                  59671,
	"openBook":                   59465,
	"sashBook":                   59672,
	"closedBook":                 61442,
	"glasses":                    59745,
	"mask":                       59777,
	"ticket":                     59788,
	"dramaMask":                  59730,
	"dice":                       59723,
	"baseball":                   59663,
	"basketball":                 59664,
	"soccerBall":                 59837,
	"tennisBall":                 59852,
	"football":                   59456,
	"lifePreserver":              59762,
	"telescope":                  59850,
	"microscope":                 59781,
	"horse":                      59756,
	"clock":                      59712,
	"alarmClock":                 59649,
	"stopwatch":                  59844,
	"bell":                       59667,
	"sparklingBell":              59838,
	"heart":                      59754,
	"star":                       59841,
	"trophy":                     59860,
	"lightbulb":                  59763,
	"lightningBolt":              59764,
	"flag":                       59736,
	"tag":                        59848,
	"key":                        59760,
	"hourglass":                  59757,
	"lock":                       59770,
	"unlockButton":               59862,
	"battery":                    59489,
	"magicWand":                  59511,
	"magicWandAlt":               59771,
	"paintbrush":                 59793,
	"pencil":                     59798,
	"paperclip":                  59794,
	"scissors":                   59824,
	"magnifyingGlass":            59772,
	"chainlink":                  59685,
	"eyedropper":                 59716,
	"hammer":                     59748,
	"wrench":                     59870,
	"tools":                      59749,
	"gear":                       59743,
	"hammerAlt":                  59473,
	"screwdriver":                59825,
	"hand":                       59751,
	"trashcan":                   59859,
	"waterDrop":                  59866,
	"mug":                        59789,
	"steamingBowl":               59842,
	"apple":                      59740,
	"carrot":                     59683,
	"fish":                       59735,
	"cake":                       59679,
	"wineBottle":                 59868,
	"martini":                    59776,
	"clothesHanger":              59713,
	"laundryMachine":             59761,
	"oven":                       59792,
	"shirt":                      59827,
	"bathtub":                    59665,
	"shower":                     59829,
	"pill":                       59461,
	"medicine":                   59815,
	"medicineBottle":             59778,
	"bandage":                    59660,
	"inhaler":                    59759,
	"stethoscope":                59843,
	"syringe":                    59847,
	"atom":                       59657,
	"chemical":                   59686,
	"cat":                        59684,
	"dog":                        59728,
	"pawPrint":                   59796,
	"thumbsUp":                   59857,
	"graduate":                   59746,
	"gift":                       59744,
	"alien":                      59651,
	"bed":                        59666,
	"stairs":                     59840,
	"rocket":                     59822,
	"map":                        61444,
	"gauge":                      61452,
	"speedometer":                61453,
	"barometer":                  61454,
	"network":                    61455,
	"rectangleStack":             61456,
	"squareStack":                61457,
	"threeDSquareStack":          61458,
	"photoStack":                 61459,
	"photoStackAlt":              61460,
	"aperture":                   61461,
	"note":                       61464,
	"noteText":                   61465,
	"noteTextPlus":               61466,
	"sendMessage":                61467,
	"addMessage":                 61468,
	"earPods":                    61475,
	"airPods":                    61476,
	"airPodsPro":                 61477,
	"radio":                      61480,
	"appleTV":                    61482,
	"homePod":                    61483,
	"appleWatchWaves":            61484,
	"iPhone":                     61486,
	"iPhoneWave":                 61487,
	"iPhoneApps":                 61488,
	"iPad":                       61489,
	"iPadAlt":                    61490,
	"iPod":                       61491,
	"babyGirl":                   59658,
	"babyBoy":                    59659,
	"child":                      59687,
	"man":                        59775,
	"woman":                      59869,
	"wheelchair":                 59806,
	"person":                     59801,
	"people2":                    59800,
	"people3":                    59799,
	"person2":                    59437,
	"personAlter":                59802,
	"personSpeech":               59804,
	"personDancer":               59803,
	"personLifting":              59807,
	"personSkiing":               59809,
	"personSnowboarding":         59810,
	"personSwimming":             59811,
	"personHiking":               59805,
	"personWalking":              59812,
	"personWalkingCane":          59813,
	"personRunning":              59808,
	"personRunningCircle":        61493,
	"personSprinting":            61494,
	"personClose":                61495,
	"personOpen":                 61496,
	"shortcuts":                  61440,
	"alertTriangle":              59650,
	"arrowCurvedLeft":            59654,
	"arrowCurvedRight":           59655,
	"bookmark":                   59670,
	"barcode":                    59661,
	"QRCode":                     59819,
	"play":                       59508,
	"boxFilled":                  59673,
	"boxOutline":                 59674,
	"braille":                    59675,
	"circleLeftArrow":            59696,
	"circleRightArrow":           59705,
	"downloadArrow":              59693,
	"circledUpArrow":             59707,
	"circledDownArrow":           59692,
	"uploadArrow":                59708,
	"circledPlay":                59699,
	"circledRewind":              59704,
	"circledPower":               59702,
	"circledStop":                59706,
	"circledFastForward":         59695,
	"circledQuestionMark":        59703,
	"circledCheckmark":           59690,
	"circledPlus":                59700,
	"circledX":                   59709,
	"circledPi":                  59698,
	"circledI":                   59697,
	"smileyFace":                 59834,
	"document":                   59725,
	"dollarSign":                 59395,
	"poundSign":                  59512,
	"euroSign":                   59448,
	"yenSign":                    59514,
	"bitcoin":                    59515,
	"asterisk":                   59656,
	"documentFilled":             59726,
	"documentOutline":            59727,
	"newsArticle":                59791,
	"fourSquares":                59739,
	"ellipsis":                   59392,
	"list":                       59445,
	"twelveSquares":              59405,
	"Connected":                  59718,
	"infinity":                   59758,
	"recycle":                    59820,
	"loadingIndicator":           59767,
	"loadingIndicatorAlt":        59516,
	"Target":                     59849,
	"podcasts":                   59816,
	"targetAlt":                  59454,
	"locationArrow":              59768,
	"locationPin":                59769,
	"parking":                    59795,
	"crop":                       59720,
	"shrinkArrow":                59830,
	"moveArrow":                  59786,
	"repostArrows":               59821,
	"syncArrows":                 59846,
	"shuffleArrows":              59832,
	"sliders":                    59833,
	"doubleQuote":                59729,
	"peaceSign":                  59797,
	"threeCircles":               59856,
	"textSymbol":                 59853,
	"feedRight":                  59732,
	"feed":                       59497,
	"wifi":                       59867,
	"airdrop":                    61501,
	"arrowDiamond":               61497,
	"directionsRight":            61498,
	"airplayAudio":               61499,
	"airplayVideo":               61500,
	"musicNoteList":              61502,
	"musicNoteAlt":               61503,
	"musicSquareStack":           61504,
	"musicWaveForm":              61505,
	"livePlay":                   61506,
	"livePhoto":                  61507,
	"sloMo":                      61508,
	"timeLapse":                  61509,
	"calendarPlus":               61510,
	"calendarExclamation":        61511,
	"timer":                      61512,
	"timerSquare":                61513,
	"compose":                    61514,
	"duplicate":                  61515,
	"nightShift":                 61518,
	"trueTone":                   61519,
	"dialMin":                    61520,
	"dialMax":                    61521,
	"QRViewFinder":               61522,
	"cameraViewFinder":           61523,
	"walletPass":                 61524,
	"appearance":                 61525,
	"noSign":                     61528,
	"command":                    61529,
	"commandCircle":              61530,
	"commandSquare":              61531,
	"blank":                      999999,
	"bumps":                      59433,
	"stripe":                     59455,
	"facetime":                   59583,
	"circledHeart":               59542,
	"documentOutlineAlt":         59496,
	"circledA":                   59520,
	"folderGear":                 61571,
	"folderOutline":              61570,
	"takeout":                    61553,
	"starHalf":                   61579,
	"sparkles":                   61581,
	"surgicalMask":               61551,
	"bear":                       61554,
	"tiger":                      61555,
	"monkey":                     61556,
	"ram":                        61557,
	"rabbit":                     61558,
	"snake":                      61559,
	"chicken":                    61560,
	"pig":                        61561,
	"mouse":                      61562,
	"cow":                        61563,
	"dragon":                     61564,
	"retroAlien":                 61565,
	"robot":                      61566,
	"ghost":                      61567,
	"poop":                       61568,
	"skull":                      61569,
	"twoxTwoRectangles":          61572,
	"twoxTwoRectanglesOutline":   61573,
	"rectangleSplit":             61574,
	"rectangleSplitThree":        61575,
	"rectangleSplitThreeOutline": 61576,
	"sendMessageOutline":         61582,
	"brainHead":                  61532,
	"brain":                      61533,
	"faceGrinning":               61534,
	"faceSmiling":                61535,
	"faceGrinningSquint":         61536,
	"faceTears":                  61537,
	"faceRolling":                61538,
	"faceWink":                   61539,
	"faceGrimacing":              61540,
	"faceLove":                   61541,
	"faceKiss":                   61542,
	"faceHearts":                 61543,
	"faceSunglasses":             61544,
	"faceStarry":                 61545,
	"memoji":                     61546,
	"handSlash":                  61584,
	"handSlashOutline":           61585,
	"thumbsUpEmoji":              61547,
	"peace":                      61548,
	"loveGesture":                61549,
	"closedFist":                 61550,
	"xSquare":                    61589,
	"checklist":                  61587,
	"doubleQuoteOutline":         61583,
	"textBox":                    61588,
	"waveform":                   61586,
	"oneProngPuzzlePiece":        61552,
	"cartoonHeart":               61577,
	"twoCartoonHearts":           61578,
	"cloudService":               59459,
}
