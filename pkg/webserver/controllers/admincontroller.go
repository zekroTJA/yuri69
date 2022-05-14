package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/models"
)

type adminController struct {
	ct *controller.Controller
}

func NewAdminController(r *routing.RouteGroup, ct *controller.Controller) {
	t := adminController{ct: ct}
	r.Get("", t.getAdmins)
	r.Put("/<id>", t.putAdmin)
	r.Delete("/<id>", t.deleteAdmin)
	r.Get("/is", t.isAdmin)
	return
}

func (t *adminController) getAdmins(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	admins, err := t.ct.GetAdmins(userid)
	if err != nil {
		return err
	}

	return ctx.Write(admins)
}

func (t *adminController) isAdmin(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.CheckAdmin(userid)
	if err != nil {
		return err
	}

	return ctx.Write(models.StatusOK)
}

func (t *adminController) putAdmin(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	adminid := ctx.Param("id")

	if adminid == "" {
		return errs.WrapUserError("invalid admin user id")
	}

	if err := t.ct.SetAdmin(userid, adminid); err != nil {
		return err
	}

	return ctx.Write(models.StatusOK)
}

func (t *adminController) deleteAdmin(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	adminid := ctx.Param("id")

	if adminid == "" {
		return errs.WrapUserError("invalid admin user id")
	}

	if err := t.ct.RemoveAdmin(userid, adminid); err != nil {
		return err
	}

	return ctx.Write(models.StatusOK)
}
