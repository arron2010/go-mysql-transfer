package global

import "sync"

type RuleMap struct {
	lockOfRuleInsMap sync.RWMutex
	ruleInsMap       map[string]*Rule
}

func NewRuleMap() *RuleMap {
	r := &RuleMap{}
	r.ruleInsMap = make(map[string]*Rule, 8)
	return r
}

func (rule *RuleMap) AddRuleIns(ruleKey string, r *Rule) {
	rule.lockOfRuleInsMap.Lock()
	defer rule.lockOfRuleInsMap.Unlock()
	rule.ruleInsMap[ruleKey] = r
}

func (rule *RuleMap) RuleIns(ruleKey string) (*Rule, bool) {
	rule.lockOfRuleInsMap.RLock()
	defer rule.lockOfRuleInsMap.RUnlock()

	r, ok := rule.ruleInsMap[ruleKey]

	return r, ok
}

func (rule *RuleMap) RuleInsExist(ruleKey string) bool {
	rule.lockOfRuleInsMap.RLock()
	defer rule.lockOfRuleInsMap.RUnlock()

	_, ok := rule.ruleInsMap[ruleKey]

	return ok
}

func (rule *RuleMap) RuleInsTotal() int {
	rule.lockOfRuleInsMap.RLock()
	defer rule.lockOfRuleInsMap.RUnlock()

	return len(rule.ruleInsMap)
}

func (rule *RuleMap) RuleInsList() []*Rule {
	rule.lockOfRuleInsMap.RLock()
	defer rule.lockOfRuleInsMap.RUnlock()

	list := make([]*Rule, 0, len(rule.ruleInsMap))
	for _, rule := range rule.ruleInsMap {
		list = append(list, rule)
	}

	return list
}

func (rule *RuleMap) RuleKeyList() []string {
	rule.lockOfRuleInsMap.RLock()
	defer rule.lockOfRuleInsMap.RUnlock()

	list := make([]string, 0, len(rule.ruleInsMap))
	for k, _ := range rule.ruleInsMap {
		list = append(list, k)
	}

	return list
}
