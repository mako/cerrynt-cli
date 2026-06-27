# **🧭 cerrynt-cli — Revised Post-MVP Roadmap (Reality-Aligned)**

This replaces the earlier Phase 2–6 roadmap.  
 It reflects the current codebase state after Phase 3 UX \+ rendering fixes.

---

# **🧱 CURRENT STATE (baseline truth)**

## **What already exists (factually true now)**

### **Architecture**

* Bubble Tea Elm-style architecture is solid  
* `app.go` owns navigation  
* screens are isolated models  
* mock data exists in `internal/data`  
* domain types exist in `internal/domain`

### **UX / rendering**

* article rendering bugfixes already applied (Phase 3 commit)  
* scroll logic works (with realistic data)  
* title rendering was corrected (no hidden margin bugs)  
* feed → article → articleview navigation works

### **Data reality**

* feeds are still **mock-driven**  
* unread count is **NOT part of domain anymore**  
* no persistence exists yet  
* no API integration yet

---

# **⚠️ CORE DESIGN REALITY SHIFT**

This is the key insight your codebase already made:

## **❌ Old assumption**

Feed contains UnreadCount (UI consumes it)

## **✅ New reality**

UI derives state from articles (or backend later)

👉 This is a **domain simplification**, not just a refactor.

It affects everything going forward.

---

# **🧭 NEW PHASE STRUCTURE**

We compress and realign phases to what actually matters now.

---

# **🟢 Phase 2 — Configuration \+ Domain Stabilization**

“Make the app real, remove fake assumptions”

### **Goal**

Replace hardcoded structure with real config \+ clean domain model.

---

## **Build**

### **1\. Config system (XDG)**

* `$XDG_CONFIG_HOME/cerrynt/config.yaml`  
* contains:  
  * feeds list  
  * api base url (optional placeholder)  
  * auth token (optional)

---

### **2\. Feed source becomes config-driven**

* `data.Feeds()` no longer hardcoded  
* reads from config

---

### **3\. Remove all UI-derived fake fields**

* ❌ no `UnreadCount` in `domain.Feed`  
* UI must NOT assume backend-provided aggregates

---

### **4\. Introduce read state (local only)**

* `$XDG_DATA_HOME/cerrynt/state.json`  
* stores:

   {  
    "read": {  
      "articleID": true  
    }  
  }

---

## **Architectural impact**

* `main.go` loads config once  
* `app.New(config)`  
* domain becomes “clean input model only”

---

## **Result**

At end of Phase 2:

✔ app runs on config  
 ✔ feed list is dynamic  
 ✔ read state persists  
 ✔ domain is simplified

---

# **🟡 Phase 3 — Async \+ Store Boundary (REAL VERSION)**

“Introduce concurrency without breaking UX correctness”

This replaces the earlier “clean async phase” definition.

---

## **Key correction vs old roadmap**

We do NOT introduce abstract async first.

We first stabilize:

domain \+ config \+ UI expectations

Then async.

---

## **Build**

### **1\. Define Store interface (minimal, correct)**

type Store interface {  
   Feeds(ctx context.Context) (\[\]domain.Feed, error)  
   Articles(ctx context.Context, feedID string) (\[\]domain.Article, error)  
}

No unread counts. No extras.

---

### **2\. Introduce async in ONE place first**

Start with:

* feed list loading async  
* spinner state

NOT everything at once.

---

### **3\. Add UI states**

Each screen:

* `loading bool`  
* `err error`

---

### **4\. App-level message handling**

* `FeedsLoadedMsg`  
* `FetchErrorMsg`

---

## **Important correction**

👉 We do NOT async everything yet.

Only:

* feed list load

Articles can stay sync initially if needed.

---

## **Result**

✔ real async exists  
 ✔ but system is still debuggable  
 ✔ no race condition explosion yet

---

# **🔵 Phase 4 — API Preparation (NOT RSS-FIRST)**

“Prepare backend integration shape, not transport logic”

This is where the old roadmap was slightly wrong.

---

## **Build**

### **1\. API client skeleton**

* implements `Store`  
* returns mocked responses first

---

### **2\. Contract alignment layer**

* domain ↔ API mapping functions

---

### **3\. Auth token wiring**

* config → API client

---

## **Key shift**

❌ RSS is NOT a core phase anymore  
 ✔ RSS becomes optional adapter

Because:

backend is already the intended source of truth

---

# **🟣 Phase 5 — Backend Integration (Rails)**

“Replace Store implementation, nothing else changes”

---

## **Build**

* `internal/api` implements `Store`  
* real HTTP calls  
* read sync:  
  * mark read endpoint  
* fetch feeds/articles

---

## **Key architectural win**

Because of Phase 3:

✔ screens do NOT change  
 ✔ app.go does NOT change  
 ✔ only main.go swaps dependency

---

# **🟠 Phase 6 — UX Completion Layer**

“Make it feel like a real product”

---

## **Build**

* article rendering polish (glamour optional)  
* browser open (`o`)  
* help overlay  
* better empty states  
* optional search/filter

---

## **Important**

This phase is now PURE UX polish.

No architectural decisions here.

---

# **🔥 KEY DIFFERENCES vs old roadmap**

## **1\. RSS is downgraded**

Old:

RSS \= Phase 4 core

New:

RSS \= optional adapter, not core path

---

## **2\. UnreadCount is fully removed from system design**

Not just deleted — conceptually removed.

---

## **3\. Async is narrower first**

Old:

everything async in Phase 3

New:

only feed loading first

---

## **4\. Store is stricter**

No “future-proof extras” yet.

Keep it minimal:

feeds \+ articles only

---

## **🧭 FINAL SUMMARY**

This is the actual evolution path now:

Phase 1 (done)  
 MVP UI \+ navigation

Phase 2  
 Config \+ persistence \+ domain cleanup

Phase 3  
 Minimal async \+ Store boundary

Phase 4  
 API contract \+ client skeleton

Phase 5  
 Rails backend integration

Phase 6  
 UX polish  
