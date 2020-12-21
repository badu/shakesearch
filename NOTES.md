# How JavaScript ate my Homework

## Preface

My biggest wish right now is to have you laugh before Christmas. I'm just imagining you having a good engineering laugh, just before taking a vacation. So, please bear with this story and me, even if it was intended to be a short pitch.

I am, probably, going to be referred as the "streaming response dude", and it's fine. In the past, I've been the "byte level optimisation dude" and "to be very precise dude" before that. I'm not a proud person - just accepting myself who I am - so, I don't mind that.

## How it all started

I've encountered the challenge you've put together on Reddit. I liked it, because it has originality.
I am so "original" as the rest of mankind : I steal ideas around the Internet, without checking what I copy-paste from Stack Overflow. For a while, I've tried to discipline myself applying the question "what problem does it solves?" all over the place. That - in the end - made me lazy, because seems most of the problems are unsolvable (in a "forever" manner) in the lifespan that we have.

## Day one

From the beginning, I became that I had no idea about suffix arrays, mostly the fact that it was in the standard Go packages. Of course - me - the one who had ripped apart some standard packages, because I had to understand how things work.

I had doubts that it was a good solution: who in the right mind doesn't write their own searching an algorithm from ground up, introducing all possible bugs due to corner cases. So, I've started googling :
- "the fastest string search" - it turns out that there is an old universal problem, which some mathemagicians tried to solve. After diagonally reading some of the PHD thesis out there, I understood I'm facing a joke : like most things on this planet, we don't have a better solution nowadays, than the one discovered in the 80s;
- "better than Boyer Moore" was next. It seems that you cannot be better than them! You spend too much time watching Netflix, instead of improving your math;
- the "academician" in me, wrote a few benchmarks (and I've deleted them in shame), with copy-paste of some implementations from Github against the suffix arrays solution. I got it all wrong it seems : you cannot beat suffix array.
  This is the point where I've realized that the problem I'm trying to solve is fake, and the question is wrong. Assuming you intended to trade RAM usage for CPU speed, the suffix array seems to do the job.
  
If you are searching for an algorithm better than Boyer Moore, then you are asking for a mixed answer. Mixed answers are not good for the health : one finds itself waking up in the middle of the night thinking about the choice it made and not being able to go back to sleep.

So, I've stopped reinventing hot water, because searching strings is a bigger problem that I will not have the time to solve. Not in this life. Worth mentioning though, I've even attempted a google term "search string assembly" - took a good look, understood some of the registry loading instructions and regretted that it is confirmed : I have so much to learn yet and I've spent time learning other things.

The taker for me, is that suffix trees have a high ram memory cost an
d are slower to index text (create the tree data structure) and typically the suffix trees are used for DNA indexing and web search engines optimization. Boyer Moore (and maybe together with Horspool) is used everywhere, for example convention programs (text search function) and widely in web search engines. Note to self : "See, Badu ? You've learned something from it, so the effort is paying!"

## Day two

I don't know how other developers are doing, but me - as a "professional" developer - I'm firstly addressing my "fears". Since your challenge contained a frontend part, I had to jump on it immediately. Not knowing React, reminded me that I'm at the age where I trade readability to highly opinionated frameworks. Not saying that VueJS is not highly opinionated, just saying that it is more readable to me. After jumping into Go language, a paradigm change happen to me: instead of looking at the fastest and early optimized solutions, I'm just looking for readable code.

As the Universe arranged things, I've just finished an Udemy VueJS course, by Maximilian Schwarzmüller. A side note: that guy is good. I'm rarely impressed by a teacher's skills, because there are so many imposters out there - as we all are aware.

Of course, as soon as I've started writing the VueJS part, I realized that "I know nothing", just like that dog said in Man In Black. 

## Day three

This day was three days apart from the previous two. Not because I gave up, but because I'm practicing 'idea distancing', in these times of 'social distancing'. I know myself really well, and I know that my imagination is not my best skill.

So, I've said to my self : it's time for some cheating : by looking at the forks that has been made to the repository, one can find what others are building for this challenge. Looking in the other people notebook might trigger an interesting idea that worth taking to the next level. It is still cheating, but hey, this is the destiny of the dudes that were born like me - a man of few ideas.

I was immediately disappointed : no-one seemed to tackle the panic triggered by taking 250 runes from "complete works" string! I mean, obviously, nobody bothered to check what happens if you are outside the string's boundary ? Really! Man, in what world we're living !?

Other than that, I've found NO cool ideas: my fellow competitors, writing functionality that exists from scratch, over engineering a challenge that was supposed to be fun and showing how smart they are. No pun intended, just an observation.

This brings me to the plot of the story. I thought : how about doing that "json response streaming" that I've planned to experiment on, for so long? "I will see what 'cool' feature I can add later on, but for the moment, seems a good idea to jump right into solving this problem" was my thought.

I've fixed the panic, fixed the sensitive search and closed the day.

## The Saga

Ok, speeding this up : I've used echo framework, despite the fact that I'm totally capable of writing all those middlewares and wiring myself.
I've used a folder structure that I like, which is highly opinionated.
I've used async-await in the Repository implementation, since there are so many questions on Reddit about how you can achieve that in Go.

pipeThrough

```json
// eslint-disable-next-line no-undef
TransformStream
```
