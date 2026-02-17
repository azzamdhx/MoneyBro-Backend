CREATE TABLE savings_contributions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    savings_goal_id UUID NOT NULL REFERENCES savings_goals(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    contribution_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_savings_contributions_goal_id ON savings_contributions(savings_goal_id);
