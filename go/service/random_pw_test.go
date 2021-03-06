package service

import (
	"testing"

	"github.com/keybase/client/go/protocol/keybase1"
	"golang.org/x/net/context"

	"github.com/keybase/client/go/kbtest"
	"github.com/keybase/client/go/libkb"
	"github.com/stretchr/testify/require"
)

func sessCheck(t *testing.T, g *libkb.GlobalContext) (err error) {
	arg := libkb.NewRetryAPIArg("sesscheck")
	arg.SessionType = libkb.APISessionTypeREQUIRED
	_, err = g.API.Get(arg)
	t.Logf("sesscheck returned: %q", err)
	return err
}

func TestSignupRandomPWUser(t *testing.T) {
	tc := libkb.SetupTest(t, "randompw", 3)
	defer tc.Cleanup()

	_, err := kbtest.CreateAndSignupFakeUserRandomPW("rpw", tc.G)
	require.NoError(t, err)

	userHandler := NewUserHandler(nil, tc.G, nil)
	ret, err := userHandler.LoadHasRandomPw(context.Background(), keybase1.LoadHasRandomPwArg{})
	require.NoError(t, err)
	require.True(t, ret)

	// Another call to test the caching
	ret, err = userHandler.LoadHasRandomPw(context.Background(), keybase1.LoadHasRandomPwArg{})
	require.NoError(t, err)
	require.True(t, ret)

	// Another one with ForceRepoll
	ret, err = userHandler.LoadHasRandomPw(context.Background(), keybase1.LoadHasRandomPwArg{ForceRepoll: true})
	require.NoError(t, err)
	require.True(t, ret)

	// TODO: Bring back some kind "you shouldn't logout" test in CORE-10084.

	handler := NewLoginHandler(nil, tc.G)
	err = handler.Logout(context.Background(), keybase1.LogoutArg{})
	require.NoError(t, err)
	require.Error(t, sessCheck(t, tc.G), "expecting to be logged out and error in sesscheck")

	// handler := NewLoginHandler(nil, tc.G)
	// err = handler.Logout(context.Background(), keybase1.LogoutArg{})
	// require.Error(t, err)
	// require.Contains(t, err.Error(), "Cannot logout")
	// require.Contains(t, err.Error(), "Set a password first")
	// require.NoError(t, sessCheck(t, tc.G), "expecting to still have a valid session and no error in sesscheck")

	// err = handler.Logout(context.Background(), keybase1.LogoutArg{Force: true})
	// require.NoError(t, err)
	// require.Error(t, sessCheck(t, tc.G), "expecting to be logged out and error in sesscheck")
}
