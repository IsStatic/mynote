package main

import (
	"context"
	"mynote/cmd/note/pack"
	"mynote/cmd/note/service"
	"mynote/idl/notedemo"
	"mynote/pkg/constants"
	"mynote/pkg/errno"
)

// NoteServiceImpl implements the last service interface defined in the IDL.
type NoteServiceImpl struct{}

// CreateNote implements the NoteServiceImpl interface.
func (s *NoteServiceImpl) CreateNote(ctx context.Context, req *notedemo.CreateNoteRequest) (resp *notedemo.CreateNoteResponse, err error) {
	resp = new(notedemo.CreateNoteResponse)

	if req.UserId <= 0 || len(req.Title) == 0 || len(req.Content) == 0 {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = service.NewCreateNoteService(ctx).CreateNote(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// MGetNote implements the NoteServiceImpl interface.
func (s *NoteServiceImpl) MGetNote(ctx context.Context, req *notedemo.MGetNoteRequest) (resp *notedemo.MGetNoteResponse, err error) {
	resp = new(notedemo.MGetNoteResponse)

	if len(req.NoteIds) == 0 {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	notes, err := service.NewMGetNoteService(ctx).MGetNote(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.Notes = notes
	return resp, nil
}

// DeleteNote implements the NoteServiceImpl interface.
func (s *NoteServiceImpl) DeleteNote(ctx context.Context, req *notedemo.DeleteNoteRequest) (resp *notedemo.DeleteNoteResponse, err error) {
	resp = new(notedemo.DeleteNoteResponse)

	if req.NoteId <= 0 {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = service.NewDelNoteService(ctx).DelNote(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// QueryNote implements the NoteServiceImpl interface.
func (s *NoteServiceImpl) QueryNote(ctx context.Context, req *notedemo.QueryNoteRequest) (resp *notedemo.QueryNoteResponse, err error) {
	resp = new(notedemo.QueryNoteResponse)

	if req.UserId <= 0 || req.Limit < 0 || req.Offset < 0 {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	if req.Limit == 0 {
		req.Limit = constants.DefaultLimit
	}

	notes, total, err := service.NewQueryNoteService(ctx).QueryNoteService(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.Notes = notes
	resp.Total = total
	return resp, nil
}

// UpdateNote implements the NoteServiceImpl interface.
func (s *NoteServiceImpl) UpdateNote(ctx context.Context, req *notedemo.UpdateNoteRequest) (resp *notedemo.UpdateNoteResponse, err error) {
	resp = new(notedemo.UpdateNoteResponse)

	if req.NoteId <= 0 {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = service.NewUpdateNoteService(ctx).UpdateNote(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}
