package cmdcomplutil_test

import (
	"errors"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewTagAutoComplete(t *testing.T) {
	bFalse := false
	tts := []struct {
		name       string
		toComplete string
		factory    func(t *testing.T) cmdutil.Factory
		err        string
		args       cmdcompl.ValidArgs
	}{
		{
			name: "allow archived disabled",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return([]dto.Tag{
					{ID: "p0dog", Name: "Tag 0"},
					{ID: "pcat", Name: "Tag 1"},
					{ID: "catp", Name: "Tag 2"},
					{ID: "pcatp", Name: "Tag 3"},
					{ID: "p4", Name: "Tag 4"},
				}, nil)

				return f
			},
			toComplete: "cat",
			args: cmdcompl.ValidArgsMap{
				"pcat":  "Tag 1",
				"catp":  "Tag 2",
				"pcatp": "Tag 3",
			},
		},
		{
			name: "allow archived enabled",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowArchivedTags: true,
				})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Tag{
					{ID: "t1", Name: "Tag 1"},
				}, nil)

				return f
			},
			args: cmdcompl.ValidArgsMap{
				"t1": "Tag 1",
			},
		},
		{
			name: "no workspace, nothing",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().
					Return("", errors.New("no workspace"))
				return f
			},
			err:  "no workspace",
			args: cmdcompl.EmptyValidArgs(),
		},
	}

	for i := range tts {
		tt := tts[i]
		t.Run(tt.name, func(t *testing.T) {
			f := tt.factory(t)
			autoComplete := cmdcomplutil.NewTagAutoComplete(f, f.Config())

			args, err := autoComplete(
				&cobra.Command{}, []string{}, tt.toComplete)

			if tt.err == "" && !assert.NoError(t, err) {
				return
			}

			if tt.err != "" && (!assert.Error(t, err) || !assert.Regexp(
				t, tt.err, err.Error())) {
				return
			}

			assert.Equal(t, tt.args, args)
		})
	}
}
