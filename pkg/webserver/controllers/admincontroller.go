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
	r.Get("/guilds", t.getGuilds)
	r.Delete("/guilds/<id>", t.removeGuilds)
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

	user, err := t.ct.SetAdmin(userid, adminid)
	if err != nil {
		return err
	}

	return ctx.Write(user)
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

func (t *adminController) getGuilds(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	guilds, err := t.ct.GetGuilds(userid)
	if err != nil {
		return err
	}

	guildsResponse := make([]models.Guild, 0, len(guilds))
	for _, g := range guilds {
		guildsResponse = append(guildsResponse, models.GuildFromGuild(g))
	}

	return ctx.Write(guildsResponse)
}

func (t *adminController) removeGuilds(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	guildid := ctx.Param("id")

	if err := t.ct.RemoveGuild(userid, guildid); err != nil {
		return err
	}

	return ctx.Write(models.StatusOK)
}
