# AI Code Review — Logger System

## Context for the AI
You are reviewing a solution to "Logger System" — an LLD interview problem testing Strategy, Observer, Factory, and Singleton patterns with HashMap + String Parsing + Queue as the DSA core. This was asked at Amazon.

The candidate was given a multi-part problem:
- Part 1: Singleton logger with level filtering (DEBUG/INFO/WARN/ERROR/FATAL), history tracking, and console output
- Part 2: Pluggable output formats via Strategy — PlainText, JSON, CSV formatters; setFormatter() at runtime
- Part 3: Multiple simultaneous output destinations via Observer — each destination has its own formatter; fault isolation between destinations

## Review Criteria

### 1. Pattern Correctness
- Is Logger a true Singleton — only one instance, private constructor, deleted copy and assignment?
- Is LogFormatter a clean Strategy interface with a single format() method?
- Is LogDestination a clean Observer interface — does Logger iterate destinations without knowing their concrete types?
- Does each destination own its own formatter (independent formatting per destination)?

### 2. Open/Closed Principle
- Can a new formatter (e.g., XmlFormatter) be added without modifying Logger?
- Can a new destination (e.g., NetworkDestination) be registered without changing Logger::log()?
- Does Logger::log() contain zero if-else logic on formatter or destination type?

### 3. C++ Quality
- Meyers Singleton (`static Logger inst;`) vs. heap-allocated (`static Logger* inst`): which is used and why?
- Fault isolation: is each dest->write() wrapped in try-catch to prevent one failure from blocking others?
- levelToString: static unordered_map for O(1) lookup vs. switch statement — which is more extensible?

### 4. Extension Handling
- Did the Part 1 formatter field naturally upgrade to a pluggable strategy in Part 2?
- Did the destinations vector in Part 3 cleanly compose with the existing history + formatter mechanism?
- What had to change at each extension point? What survived unchanged?

### 5. Interview Readiness
- Could the candidate explain why a Singleton is appropriate here and when it would be a design smell?
- What follow-up questions would expose weak understanding?
  - "How would you make this thread-safe for a multi-threaded server?"
  - "How would you support async (non-blocking) log writes to a file destination?"
  - "If two formatters exist (one for history, one per-destination), how does the candidate manage this?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
