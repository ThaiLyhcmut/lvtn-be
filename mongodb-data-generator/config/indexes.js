module.exports = {
  users: [
    { code: 1 },
    { email: 1 },
    { role_id: 1 },
    { major: 1 },
    { status: 1 },
    { created_at: -1 }
  ],
  
  departments: [
    { code: 1 },
    { head_id: 1 },
    { is_active: 1 }
  ],
  
  theses: [
    { title: 'text' },
    { student_id: 1 },
    { supervisor_id: 1 },
    { status_id: 1 },
    { academic_year: 1, semester: 1 },
    { created_at: -1 }
  ],
  
  thesis_statuses: [
    { name: 1 },
    { order: 1 }
  ],
  
  supervisor_assignments: [
    { thesis_id: 1 },
    { supervisor_id: 1 },
    { role: 1 },
    { is_active: 1 }
  ],
  
  submissions: [
    { thesis_id: 1 },
    { submitted_by: 1 },
    { type: 1 },
    { status: 1 },
    { submitted_at: -1 }
  ],
  
  reviews: [
    { submission_id: 1 },
    { reviewer_id: 1 },
    { status: 1 },
    { reviewed_at: -1 }
  ],
  
  defense_schedules: [
    { thesis_id: 1 },
    { defense_date: 1 },
    { status: 1 },
    { location: 1, defense_date: 1 }
  ],
  
  defense_scores: [
    { defense_schedule_id: 1 },
    { scorer_id: 1 },
    { criteria: 1 }
  ],
  
  event_logs: [
    { user_id: 1 },
    { action: 1 },
    { entity_type: 1, entity_id: 1 },
    { timestamp: -1 }
  ],
  
  archived_theses: [
    { original_thesis_id: 1 },
    { title: 'text' },
    { graduation_year: -1 },
    { archived_at: -1 }
  ],
  
  archived_submissions: [
    { archived_thesis_id: 1 },
    { original_submission_id: 1 },
    { type: 1 }
  ],
  
  archived_reviews: [
    { archived_submission_id: 1 },
    { original_review_id: 1 }
  ],

  mongoURI: "mongodb://thaily:Th%40i2004@localhost:27017/lvtn?authSource=admin"
};