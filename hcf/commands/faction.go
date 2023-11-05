package commands

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/area"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
	"time"
)

// FactionCreate is a command used to create a faction.
type FactionCreate struct {
	Sub  cmd.SubCommand `cmd:"create"`
	Name string         `cmd:"name"`
}

// FactionDisband is a command used to disband a faction.
type FactionDisband struct {
	Sub cmd.SubCommand `cmd:"disband"`
}

// FactionInformation is a command used to get information about a faction.
type FactionInformation struct {
	Sub  cmd.SubCommand            `cmd:"info"`
	Name cmd.Optional[factionName] `optional:"" cmd:"name"`
	srv  *server.Server
}

func NewFactionInformation(srv *server.Server) FactionInformation {
	return FactionInformation{srv: srv}
}

// FactionWho is a command used to get information about a faction.
type FactionWho struct {
	Sub  cmd.SubCommand              `cmd:"who"`
	Name cmd.Optional[factionMember] `optional:"" cmd:"name"`
	srv  *server.Server
}

func NewFactionWho(srv *server.Server) FactionWho {
	return FactionWho{srv: srv}
}

// FactionInvite is a command used to invite a player to a faction.
type FactionInvite struct {
	Sub    cmd.SubCommand `cmd:"invite"`
	Target []cmd.Target   `cmd:"target"`
}

// FactionJoin is a command used to join a faction.
type FactionJoin struct {
	Sub        cmd.SubCommand `cmd:"join"`
	Invitation invitation     `cmd:"invitation"`
}

// FactionLeave is a command used to leave a faction.
type FactionLeave struct {
	Sub cmd.SubCommand `cmd:"leave"`
}

// FactionKick is a command used to kick a player from a faction.
type FactionKick struct {
	Sub    cmd.SubCommand `cmd:"kick"`
	Member member         `cmd:"member"`
}

// FactionPromote is a command used to promote a player in a faction.
type FactionPromote struct {
	Sub    cmd.SubCommand `cmd:"promote"`
	Member member         `cmd:"member"`
}

// FactionDemote is a command used to demote a player in a faction.
type FactionDemote struct {
	Sub    cmd.SubCommand `cmd:"demote"`
	Member member         `cmd:"member"`
}

// FactionTop is a command used to get the top factions.
type FactionTop struct {
	Sub cmd.SubCommand `cmd:"top"`
}

// FactionClaim is a command used to claim land for a faction.
type FactionClaim struct {
	Sub cmd.SubCommand `cmd:"claim"`
}

// FactionUnClaim is a command used to unclaim land for a faction.
type FactionUnClaim struct {
	Sub cmd.SubCommand `cmd:"unclaim"`
}

// FactionSetHome is a command used to set a faction's home.
type FactionSetHome struct {
	Sub cmd.SubCommand `cmd:"sethome"`
}

// FactionHome is a command used to teleport to a faction's home.
type FactionHome struct {
	Sub cmd.SubCommand `cmd:"home"`
}

// FactionList is a command used to list factions.
type FactionList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// FactionFocusFaction is a command used to focus on a faction.
type FactionFocusFaction struct {
	Sub  cmd.SubCommand `cmd:"focus"`
	Name factionName    `cmd:"name"`
}

// FactionFocusPlayer is a command used to focus on a player.
type FactionFocusPlayer struct {
	Sub    cmd.SubCommand `cmd:"focus"`
	Target []cmd.Target   `cmd:"target"`
}

// FactionUnFocus is a command used to unfocus on a faction.
type FactionUnFocus struct {
	Sub cmd.SubCommand `cmd:"unfocus"`
}

// FactionChat is a command used to chat in a faction.
type FactionChat struct {
	Sub cmd.SubCommand `cmd:"chat"`
}

// FactionWithdraw is a command used to withdraw money from a faction.
type FactionWithdraw struct {
	Sub     cmd.SubCommand `cmd:"withdraw"`
	Balance float64        `cmd:"balance"`
}

// FactionDeposit is a command used to deposit money into a faction.
type FactionDeposit struct {
	Sub     cmd.SubCommand `cmd:"deposit"`
	Balance float64        `cmd:"balance"`
}

// FactionWithdrawAll is a command used to withdraw all the money from a faction.
type FactionWithdrawAll struct {
	Sub cmd.SubCommand `cmd:"withdraw"`
	All cmd.SubCommand `cmd:"all"`
}

// FactionDepositAll is a command used to deposit all of a user's money into a faction.
type FactionDepositAll struct {
	Sub cmd.SubCommand `cmd:"deposit"`
	All cmd.SubCommand `cmd:"all"`
}

// FactionStuck is a command to teleport to a safe position
type FactionStuck struct {
	Sub cmd.SubCommand `cmd:"stuck"`
}

// Run ...
func (f FactionCreate) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		return
	}
	if _, ok := u.FactionCreateDelay(); ok {
		o.Error(lang.Translatef(l, "command.faction.create.wait"))
		return
	}
	name := s.(cmd.NamedTarget).Name()

	if _, ok := factions.LookupMemberName(name); ok {
		o.Error(lang.Translatef(l, "command.faction.create.already.in"))
		return
	}
	if strings.TrimSpace(f.Name) == "" {
		o.Error(lang.Translatef(l, "command.faction.create.usage"))
		return
	}
	if len(f.Name) > 12 {
		o.Error(lang.Translatef(l, "command.faction.create.length.maximum"))
		return
	}
	if len(f.Name) < 3 {
		o.Error(lang.Translatef(l, "command.faction.create.length.minimum"))
		return
	}
	_, ok = factions.LookupName(f.Name)
	if ok {
		o.Error(lang.Translatef(l, "command.faction.create.exists"))
		return
	}
	fa := factions.NewFaction(f.Name, []*factions.Member{factions.NewMember(strings.ToLower(name), name, factions.RankLeader{})}, 1.1, mgl64.Vec3{}, 0, 0, 0, util.AreaVec2{})
	if err := data.SaveFaction(fa); err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
	u.SetFactionCreateDelay()
	u.SetFaction(fa)
	u.UpdateState()
	r := u.Ranks().Highest()
	_, _ = chat.Global.WriteString(lang.Translatef(l, "command.faction.create", f.Name, r.Tag(name)))
}

// Run ...
func (f FactionDisband) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	_, ok = user.Lookup(p)
	if !ok {
		return
	}
	name := p.Name()
	fa, ok := factions.LookupMemberName(name)
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if !strings.EqualFold(fa.Leader().Name(), name) {
		o.Error(lang.Translatef(l, "command.faction.disband.must.leader"))
		return
	}
	players := fa.Members()
	if err := data.DeleteFaction(fa); err != nil {
		o.Error(lang.Translatef(l, "command.faction.disband.fail"))
		return
	}
	for _, m := range players {
		mem, ok := user.LookupName(m.Name())
		if ok {
			mem.UpdateState()
			mem.SetFaction(nil)
			mem.Player().Message(lang.Translatef(l, "command.faction.disband.disbanded", name))
		}
	}
}

// Run ...
func (f FactionInformation) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	n, _ := f.Name.Load()
	name := string(n)
	if strings.TrimSpace(name) == "" {
		fa, ok := factions.LookupMemberName(sourceName)
		if !ok {
			o.Error(lang.Translatef(l, "user.faction-less"))
			return
		}
		o.Print(fa.Information(f.srv))
		return
	}
	var anyFound bool

	fa, ok := factions.LookupName(name)
	if ok {
		o.Print(fa.Information(f.srv))
		anyFound = true
	}
	fa, ok = factions.LookupMemberName(name)
	if ok {
		o.Print(fa.Information(f.srv))
		anyFound = true
	}
	if !anyFound {
		o.Error(lang.Translatef(l, "command.faction.info.not.found", name))
		return
	}
}

// Run ...
func (f FactionWho) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	n, _ := f.Name.Load()
	name := string(n)
	if strings.TrimSpace(name) == "" {
		fa, ok := factions.LookupMemberName(sourceName)
		if !ok {
			o.Error(lang.Translatef(l, "user.faction-less"))
			return
		}
		o.Print(fa.Information(f.srv))
		return
	}
	var anyFound bool

	fa, ok := factions.LookupName(name)
	if ok {
		o.Print(fa.Information(f.srv))
		anyFound = true
	}
	fa, ok = factions.LookupMemberName(name)
	if ok {
		o.Print(fa.Information(f.srv))
		anyFound = true
	}
	if !anyFound {
		o.Error(lang.Translatef(l, "command.faction.info.not.found", name))
		return
	}
}

// Run ...
func (f FactionInvite) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	src, ok := user.LookupName(s.(cmd.NamedTarget).Name())
	if !ok {
		return
	}
	fa, ok := src.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if !strings.EqualFold(fa.Leader().Name(), src.Name()) && !fa.IsCoLeader(src.Name()) && !fa.IsCaptain(src.Name()) {
		o.Error(lang.Translatef(l, "command.faction.invite.missing.permission"))
		return
	}
	if len(f.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	target := f.Target[0]
	for _, m := range fa.Members() {
		if strings.EqualFold(m.Name(), target.(cmd.NamedTarget).Name()) {
			o.Error(lang.Translatef(l, "command.faction.invite.already.in", m.DisplayName()))
			return
		}
	}
	p, ok := target.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	u.Player().Message(lang.Translatef(l, "command.faction.invite.received", fa.DisplayName(), src.Ranks().Highest().Tag(src.Name())))
	u.Invite(fa.Name())
	for _, m := range fa.Members() {
		u, ok := user.LookupName(m.Name())
		if ok {
			u.Player().Message(lang.Translatef(l, "command.faction.invite.sent", p.Name()))
		}
	}
}

// Run ...
func (f FactionJoin) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if ok {
		o.Error(lang.Translatef(l, "user.already.in.faction"))
		return
	}
	p := s.(*player.Player)
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	if !slices.Contains(u.Invitations(), string(f.Invitation)) {
		o.Error(lang.Translatef(l, "command.faction.join.not.invited"))
		return
	}
	fa, ok = factions.LookupName(strings.ToLower(string(f.Invitation)))
	if !ok {
		o.Error(lang.Translatef(l, "command.faction.join.not.found"))
		return
	}
	if fa.PlayerCount() >= 7 {
		o.Error(lang.Translatef(l, "command.faction.join.full"))
		return
	}

	fa.AddMember(p)
	err := data.SaveFaction(fa)
	if err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
	u.SetFaction(fa)
	u.UnInvite(string(f.Invitation))
	o.Print(lang.Translatef(l, "command.faction.join.joined", fa.DisplayName()))
	for _, m := range fa.Members() {
		if mem, ok := user.LookupName(m.Name()); ok {
			mem.UpdateState()
			mem.Player().Message(lang.Translatef(l, "command.faction.join.user.joined", p.Name()))
		}
	}
}

// Run ...
func (f FactionLeave) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if strings.EqualFold(fa.Leader().Name(), sourceName) {
		o.Error(lang.Translatef(l, "command.faction.leave.leader"))
		return
	}
	players := fa.Members()
	p := s.(*player.Player)
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa.RemoveMember(p)
	u.SetFaction(nil)
	for _, m := range fa.Members() {
		if mem, ok := user.LookupName(m.Name()); ok {
			mem.UpdateState()
		}
	}
	err := data.SaveFaction(fa)
	if err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
	for _, m := range players {
		if u, ok := user.LookupName(m.Name()); ok {
			u.Player().Message(lang.Translatef(l, "command.faction.leave.user.left", p.Name()))
		}
	}
}

// Run ...
func (f FactionKick) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if !strings.EqualFold(fa.Leader().Name(), sourceName) && !fa.IsCoLeader(sourceName) && !fa.IsCaptain(sourceName) {
		o.Error(lang.Translatef(l, "command.faction.kick.missing.permission"))
		return
	}
	if string(f.Member) == sourceName {
		o.Error(lang.Translatef(l, "command.faction.kick.self"))
		return
	}
	if strings.EqualFold(fa.Leader().Name(), string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.kick.leader"))
		return
	}
	if fa.IsCoLeader(sourceName) && fa.IsCoLeader(string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.kick.co-leader"))
		return
	}
	if fa.IsCaptain(sourceName) && fa.IsCaptain(string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.kick.captain"))
		return
	}
	if _, ok := fa.Member(string(f.Member)); !ok {
		o.Error(lang.Translatef(l, "command.faction.kick.not.found", string(f.Member)))
		return
	}
	u, ok := user.LookupName(string(f.Member))
	if ok {
		u.SetFaction(nil)
		u.Player().Message(lang.Translatef(l, "command.faction.kick.user.kicked"))
	}
	for _, m := range fa.Members() {
		if mem, ok := user.LookupName(m.Name()); ok {
			mem.UpdateState()
		}
	}

	fa.RemoveMemberName(string(f.Member))

	for _, m := range fa.Members() {
		if u, ok := user.LookupName(m.Name()); ok {
			u.Player().Message(lang.Translatef(l, "command.faction.kick.user.kicked", string(f.Member)))
		}
	}

	err := data.SaveFaction(fa)
	if err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
}

// Run ...
func (f FactionPromote) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if m, ok := fa.Member(sourceName); !ok || m.Rank() == (factions.RankMember{}) {
		o.Error(lang.Translatef(l, "command.faction.promote.missing.permission"))
		return
	}
	if string(f.Member) == sourceName {
		o.Error(lang.Translatef(l, "command.faction.promote.self"))
		return
	}
	if strings.EqualFold(fa.Leader().Name(), string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.promote.leader"))
		return
	}
	if fa.IsCoLeader(sourceName) && fa.IsCoLeader(string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.promote.co-leader"))
		return
	}
	if fa.IsCaptain(sourceName) && fa.IsCaptain(string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.promote.captain"))
		return
	}
	if _, ok := fa.Member(string(f.Member)); !ok {
		o.Error(lang.Translatef(l, "command.faction.member.not.found", string(f.Member)))
		return
	}
	fa.Promote(string(f.Member))
	err := data.SaveFaction(fa)
	if err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
	rankName := "Captain"
	if strings.EqualFold(fa.Leader().Name(), string(f.Member)) {
		rankName = "Leader"
	}
	for _, m := range fa.Members() {
		if u, ok := user.LookupName(m.Name()); ok {
			u.Message("command.faction.promote.user.promoted", string(f.Member), rankName)
		}
	}
}

// Run ...
func (f FactionDemote) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	sourceName := s.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}
	if !strings.EqualFold(fa.Leader().Name(), sourceName) {
		o.Error(lang.Translatef(l, "command.faction.demote.missing.permission"))
		return
	}
	if strings.EqualFold(string(f.Member), sourceName) {
		o.Error(lang.Translatef(l, "command.faction.demote.self"))
		return
	}
	if strings.EqualFold(fa.Leader().Name(), string(f.Member)) {
		o.Error(lang.Translatef(l, "command.faction.demote.leader"))
		return
	}
	if m, ok := fa.Member(string(f.Member)); !ok || m.Rank() == (factions.RankMember{}) {
		if m.Rank() == (factions.RankMember{}) {
			o.Error(lang.Translatef(l, "command.faction.demote.member"))
		} else {
			o.Error(lang.Translatef(l, "command.faction.member.not.found", string(f.Member)))
		}
		return
	}
	fa.Demote(string(f.Member))
	err := data.SaveFaction(fa)
	if err != nil {
		o.Error(lang.Translatef(l, "faction.save.fail"))
		return
	}
	for _, m := range fa.Members() {
		if u, ok := user.LookupName(m.Name()); ok {
			u.Message("command.faction.demote.user.demoted", string(f.Member), "Member")
		}
	}
}

// Run ...
func (f FactionTop) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	facs := factions.All()

	if len(facs) == 0 {
		o.Error("There are no factions.")
		return
	}

	sort.Slice(facs, func(i, j int) bool {
		return facs[i].Points() > facs[j].Points()
	})

	var top string
	top += text.Colourf("        <yellow>Top factions</yellow>\n")
	top += "\uE000\n"
	userFaction, ok := factions.LookupMemberName(u.Name())
	for i, tm := range facs {
		if i > 9 {
			break
		}
		if ok && userFaction == tm {
			top += text.Colourf(" <grey>%d. <green>%s</green> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName(), tm.Points())
		} else {
			top += text.Colourf(" <grey>%d. <red>%s</red> (<yellow>%d</yellow>)</grey>\n", i+1, tm.DisplayName(), tm.Points())
		}
	}
	top += "\uE000"
	p.Message(top)
}

// Run ...
func (f FactionClaim) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
	}
	if _, ok := fa.Claim(); ok {
		o.Error("Your faction already has a claim.")
		return
	}
	_, _ = p.Inventory().AddItem(item.NewStack(item.Hoe{Tier: item.ToolTierDiamond}, 1).WithValue("CLAIM_WAND", true))
}

// Run ...
func (f FactionUnClaim) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
	}
	if !strings.EqualFold(fa.Leader().Name(), u.Name()) {
		o.Error("You are not the faction leader.")
		return
	}
	fa.UnClaim()
	fa.SetHome(mgl64.Vec3{})
	u.Message("command.unclaim.success")
}

// Run ...
func (f FactionSetHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
	}
	c, ok := fa.Claim()
	if !ok {
		o.Error("Your faction does not have a claim.")
		return
	}
	if !c.Vec3WithinOrEqualXZ(p.Position()) {
		o.Error("You are not within your faction's claim.")
		return
	}
	if !strings.EqualFold(fa.Leader().Name(), u.Name()) && !fa.IsCaptain(u.Name()) {
		o.Error("You are not the faction leader or captain.")
		return
	}
	fa.SetHome(p.Position())
	o.Print("Your faction's home has been set.")
}

// Run ...
func (f FactionHome) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	if u.Tags().Combat().Active() {
		o.Error("You cannot teleport while in combat.")
		return
	}

	if u.Teleportations().Home().Teleporting() {
		o.Error("You are already teleporting.")
		return
	}

	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
		return
	}
	h, ok := fa.Home()
	if !ok {
		o.Error("Your faction does not have a home.")
		return
	}
	if area.Spawn(p.World()).Area().Vec3WithinOrEqualXZ(p.Position()) {
		p.Teleport(h)
		return
	}

	dur := time.Second * 10
	for t, c := range factions.Claims() {
		if t != fa.Faction && c.Vec3WithinOrEqualXZ(p.Position()) {
			dur = time.Second * 20
			break
		}
	}
	u.Teleportations().Home().Teleport(p, dur, h)
}

func onlineCount(fa *factions.Faction) int {
	var count int
	for _, m := range fa.Members() {
		_, ok := user.LookupName(m.Name())
		if ok {
			count++
		}
	}
	return count
}

// Run ...
func (f FactionList) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	facs := factions.All()
	if len(facs) == 0 {
		o.Error("There are no factions.")
		return
	}
	sort.Slice(facs, func(i, j int) bool {
		return onlineCount(facs[i]) > onlineCount(facs[j])
	})
	sort.Slice(facs, func(i, j int) bool {
		return facs[i].DTR() < facs[j].DTR()
	})

	for _, fa := range facs {
		if fa.DTR() <= 0 {
			facs = append(facs[:0], facs[1:]...)
		}
	}

	var list string
	list += text.Colourf("        <yellow>faction List</yellow>\n")
	list += "\uE000\n"
	userFaction, ok := factions.LookupMemberName(u.Name())
	for i, fa := range facs {
		if i > 9 {
			break
		}

		dtr := text.Colourf("<green>%.1f■</green>", fa.DTR())
		if fa.DTR() < 5 {
			dtr = text.Colourf("<yellow>%.1f■</yellow>", fa.DTR())
		}
		if fa.DTR() <= 0 {
			dtr = text.Colourf("<red>%.1f■</red>", fa.DTR())
		}
		if ok && userFaction == fa {
			list += text.Colourf(" <grey>%d. <green>%s</green> (<green>%d/%d</green>)</grey> %s <yellow>DTR</yellow>\n", i+1, fa.DisplayName(), onlineCount(fa), len(fa.Members()), dtr)
		} else {
			list += text.Colourf(" <grey>%d. <red>%s</red> (<green>%d/%d</green>)</grey> %s <yellow>DTR</yellow>\n", i+1, fa.DisplayName(), onlineCount(fa), len(fa.Members()), dtr)
		}
	}
	list += "\uE000"
	p.Message(list)
}

// Run ...
func (f FactionFocusFaction) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
		return
	}
	targetFaction, ok := factions.LookupName(string(f.Name))
	if !ok {
		o.Error("That faction does not exist.")
		return
	}

	if fa.Faction == targetFaction {
		o.Error("You cannot focus your own faction.")
		return
	}

	fa.FocusFaction(targetFaction)

	members, _ := u.Focusing()
	for _, m := range members {
		p, ok := user.LookupName(m.Name())
		if !ok {
			continue
		}
		p.UpdateState()
	}

	fa.Broadcast("command.faction.focus", targetFaction.DisplayName())
}

// Run ...
func (f FactionUnFocus) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
		return
	}

	var name string

	focused, okFaction := fa.FocusedFaction()
	focusedPlayer, okPlayer := fa.FocusedPlayer()

	if okFaction {
		name = focused.DisplayName()
	} else if okPlayer {
		name = focusedPlayer
	} else {
		o.Error("Your faction is not focusing another faction.")
		return
	}
	members, _ := u.Focusing()
	fa.UnFocus()

	for _, m := range members {
		p, ok := user.LookupName(m.Name())
		if !ok {
			continue
		}
		p.UpdateState()
	}

	fa.Broadcast("command.faction.unfocus", name)
}

// Run ...
func (f FactionFocusPlayer) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		o.Error("You are not in a faction.")
		return
	}
	if len(f.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	trg := f.Target[0]
	target, ok := trg.(*player.Player)
	if !ok {
		o.Error("You must target a player.")
		return
	}
	targetUser, ok := user.Lookup(target)
	if !ok {
		o.Error("That player is not online.")
		return
	}

	if targetUser == u {
		o.Error("You cannot focus yourself.")
		return
	}

	if _, ok := fa.Member(targetUser.Name()); ok {
		o.Error("You cannot focus a member of your faction.")
		return
	}

	fa.FocusPlayer(targetUser.Name())
	targetUser.UpdateState()
	fa.Broadcast("command.faction.focus", target.Name())
}

// Run ...
func (f FactionChat) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	_, ok = u.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}

	switch u.ChatType() {
	case user.ChatTypeFaction():
		u.UpdateChatType(user.ChatTypeGlobal())
	case user.ChatTypeGlobal():
		u.UpdateChatType(user.ChatTypeFaction())
	}
}

// Run ...
func (t FactionWithdraw) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	fa, ok := u.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}

	if !strings.EqualFold(fa.Leader().Name(), p.Name()) && !fa.IsCaptain(p.Name()) {
		o.Error("You cannot withdraw any balance from your faction.")
		return
	}

	bal := t.Balance
	if bal < 1 {
		o.Error("You must provide a minimum balance of $1.")
		return
	}

	if bal > fa.Balance() {
		o.Errorf("Your faction does not have a balance of $%.2f.", bal)
		return
	}

	fa.ReduceBalance(bal)
	u.IncreaseBalance(bal)

	o.Print(text.Colourf("<green>You withdrew $%.2f from %s.</green>", bal, fa.DisplayName()))
}

// Run ...
func (f FactionDeposit) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	fa, ok := u.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}

	bal := f.Balance
	if bal < 1 {
		o.Error("You must provide a minimum balance of $1.")
		return
	}

	if bal > u.Balance() {
		o.Errorf("You do not have a balance of $%.2f.", bal)
		return
	}

	u.ReduceBalance(bal)
	fa.IncreaseBalance(bal)

	o.Print(text.Colourf("<green>You deposited $%.2f into %s.</green>", bal, fa.DisplayName()))
}

// Run ...
func (f FactionWithdrawAll) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	fa, ok := u.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}

	if !strings.EqualFold(fa.Leader().Name(), p.Name()) && !fa.IsCaptain(p.Name()) {
		o.Error("You cannot withdraw any balance from your faction.")
		return
	}

	bal := fa.Balance()
	if bal < 1 {
		o.Error("Your faction's balance is lower than $1.")
		return
	}

	fa.ReduceBalance(bal)
	u.IncreaseBalance(bal)

	o.Print(text.Colourf("<green>You withdrew $%d from %s.</green>", bal, fa.Name()))
}

// Run ...
func (f FactionDepositAll) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	fa, ok := u.Faction()
	if !ok {
		o.Error(lang.Translatef(l, "user.faction-less"))
		return
	}

	bal := u.Balance()
	if bal < 1 {
		o.Error("Your balance is lower than $1.")
		return
	}

	u.ReduceBalance(bal)
	fa.IncreaseBalance(bal)

	o.Print(text.Colourf("<green>You deposited $%d into %s.</green>", bal, fa.Name()))
}

// Run ...
func (f FactionStuck) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	pos := safePosition(u, 24)
	if pos == (cube.Pos{}) {
		u.Message("command.faction.stuck.no-safe")
		return
	}

	if u.Teleportations().Logout().Teleporting() {
		o.Error("You are already stucking.")
		return
	}

	u.Message("command.faction.stuck.teleporting", pos.X(), pos.Y(), pos.Z(), 30)
	u.Teleportations().Stuck().Teleport(p, time.Second*30, mgl64.Vec3{
		float64(pos.X()),
		float64(pos.Y()),
		float64(pos.Z()),
	})
}

func safePosition(u *user.User, radius int) cube.Pos {
	pos := cube.PosFromVec3(u.Position())
	minX := pos.X() - radius
	maxX := pos.X() + radius
	minZ := pos.Z() - radius
	maxZ := pos.Z() + radius

	for x := minX; x < maxX; x++ {
		for z := minZ; z < maxZ; z++ {
			at := pos.Add(cube.Pos{x, 0, z})
			for tm, claim := range factions.Claims() {
				if claim.Vec3WithinOrEqualXZ(at.Vec3Centre()) {
					if f, ok := u.Faction(); ok && f.Compare(tm) {
						y := u.World().Range().Max()
						for y > pos.Y() {
							y--
							b := u.World().Block(cube.Pos{x, y, z})
							if b != (block.Air{}) {
								return cube.Pos{x, y, z}
							}
						}
					}
				}
			}

			for _, ar := range append(area.Protected(u.World()), area.Wilderness(u.World())) {
				if ar.Area().Vec3WithinOrEqualXZ(at.Vec3Centre()) {
					y := u.World().Range().Max()
					for y > pos.Y() {
						y--
						b := u.World().Block(cube.Pos{x, y, z})
						if b != (block.Air{}) {
							return cube.Pos{x, y, z}
						}
					}
				}
			}
		}
	}
	return cube.Pos{}
}

// Allow ...
func (FactionCreate) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionDisband) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionInformation) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionWho) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionInvite) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionJoin) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionLeave) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionKick) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionPromote) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionDemote) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionTop) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionClaim) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionUnClaim) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionSetHome) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionHome) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionList) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionFocusFaction) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionUnFocus) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionFocusPlayer) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionChat) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionWithdraw) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionDeposit) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionWithdrawAll) Allow(s cmd.Source) bool {
	return allow(s, false)
}

// Allow ...
func (FactionDepositAll) Allow(s cmd.Source) bool {
	return allow(s, false)
}

type (
	invitation    string
	member        string
	factionName   string
	factionMember string
)

// Type ...
func (member) Type() string {
	return "member"
}

// Type ...
func (invitation) Type() string {
	return "invitation"
}

// Type ...
func (factionName) Type() string {
	return "faction_name"
}

// Type ...
func (factionMember) Type() string {
	return "faction_member"
}

// Options ...
func (invitation) Options(src cmd.Source) []string {
	p, ok := src.(*player.Player)
	if !ok {
		return nil
	}
	u, ok := user.Lookup(p)
	if !ok {
		return nil
	}
	return u.Invitations()
}

// Options ...
func (member) Options(src cmd.Source) []string {
	sourceName := src.(cmd.NamedTarget).Name()
	fa, ok := factions.LookupMemberName(sourceName)
	if !ok {
		return nil
	}
	var members []string
	for _, m := range fa.Members() {
		if !strings.EqualFold(m.Name(), sourceName) {
			members = append(members, m.DisplayName())
		}
	}
	return members
}

// Options ...
func (factionName) Options(cmd.Source) []string {
	var facs []string
	for _, tm := range factions.All() {
		facs = append(facs, tm.DisplayName())
	}
	return facs
}

// Options ...
func (factionMember) Options(cmd.Source) []string {
	var members []string
	for _, tm := range factions.All() {
		for _, m := range tm.Members() {
			members = append(members, m.DisplayName())
		}
	}
	return members
}
