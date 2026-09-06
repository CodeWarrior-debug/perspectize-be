package postgres

import "github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"

func messageModelToDomain(m *MessageModel) domain.Message {
	return domain.Message{
		ID: m.ID, ThreadID: int(m.ThreadID), SenderID: int(m.SenderID),
		Seq: m.Seq, Body: m.Body, ClientNonce: m.ClientNonce, CreatedAt: m.CreatedAt,
	}
}

func messageDomainToModel(d *domain.Message) *MessageModel {
	return &MessageModel{
		ID: d.ID, ThreadID: int64(d.ThreadID), SenderID: int64(d.SenderID),
		Body: d.Body, ClientNonce: d.ClientNonce,
	}
}

func threadParticipantModelToDomain(m *ThreadParticipantModel) domain.ThreadParticipant {
	return domain.ThreadParticipant{
		ThreadID: int(m.ThreadID), UserID: int(m.UserID),
		Role: domain.ThreadRole(m.Role), LastReadSeq: m.LastReadSeq,
		Muted: m.Muted, JoinedAt: m.JoinedAt, LeftAt: m.LeftAt,
	}
}

func messageThreadModelToDomain(m *MessageThreadModel, parts []ThreadParticipantModel) domain.MessageThread {
	t := domain.MessageThread{
		ID: int(m.ID), Title: m.Title, CreatedBy: int(m.CreatedBy),
		LastMessageAt: m.LastMessageAt, CreatedAt: m.CreatedAt,
	}
	for i := range parts {
		t.Participants = append(t.Participants, threadParticipantModelToDomain(&parts[i]))
	}
	return t
}
