//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package v1

import (
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/db"
	"github.com/cloustone/sentel/apiserver/util"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Rule Api
type ruleAddRequest struct {
	requestBase
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Target      string `json:"target"`
	ProductId   string `json:"productId"`
}

// addRule add new rule for product
func addRule(ctx echo.Context) error {
	req := new(ruleAddRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := db.Rule{
		Id:           uuid.NewV4().String(),
		Name:         req.Name,
		ProductId:    req.ProductId,
		Method:       req.Method,
		Target:       req.Target,
		TimeCreated:  time.Now(),
		TimeModified: time.Now(),
	}
	if err := r.RegisterRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"todo",
		util.TopicNameRule,
		&util.RuleTopic{
			RuleId:    rule.Id,
			ProductId: rule.ProductId,
			Action:    util.ObjectActionRegister,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}

// deleteRule delete existed rule
func deleteRule(ctx echo.Context) error {
	id := ctx.Param("id")
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if err := r.DeleteRule(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"todo",
		util.TopicNameRule,
		&util.RuleTopic{
			RuleId: id,
			Action: util.ObjectActionDelete,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String()})
}

// UpdateRule update existed rule
func updateRule(ctx echo.Context) error {
	req := new(ruleAddRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := db.Rule{
		Id:           req.Id,
		Name:         req.Name,
		ProductId:    req.ProductId,
		Method:       req.Method,
		Target:       req.Target,
		TimeModified: time.Now(),
	}
	if err := r.UpdateRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"todo",
		util.TopicNameRule,
		&util.RuleTopic{
			RuleId:    rule.Id,
			ProductId: rule.ProductId,
			Action:    util.ObjectActionUpdate,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}

// getRule retrieve a rule
func getRule(ctx echo.Context) error {
	id := ctx.Param("id")
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule, err := r.GetRule(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}
