const { ObjectId } = require('mongodb');
const faker = require('faker');
const fs = require('fs');
const path = require('path');
const { program } = require('commander');
const chalk = require('chalk');

faker.locale = 'vi';

program
  .option('-u, --users <number>', 'Number of users', '100')
  .option('-t, --theses <number>', 'Number of theses', '50')
  .option('-s, --submissions <number>', 'Number of submissions', '150')
  .option('-r, --reviews <number>', 'Number of reviews', '150')
  .option('-d, --defenses <number>', 'Number of defense schedules', '30')
  .option('-a, --archived <number>', 'Number of archived theses', '20')
  .parse(process.argv);

const options = program.opts();

const counts = {
  users: parseInt(options.users),
  theses: parseInt(options.theses),
  submissions: parseInt(options.submissions),
  reviews: parseInt(options.reviews),
  defenses: parseInt(options.defenses),
  archived: parseInt(options.archived)
};

console.log(chalk.blue('MongoDB Data Generator'));
console.log(chalk.gray('======================'));
console.log(chalk.yellow('Generating data with following counts:'));
Object.entries(counts).forEach(([key, value]) => {
  console.log(chalk.green(`  ${key}: ${value}`));
});

// Helper functions
const randomElement = (array) => array[Math.floor(Math.random() * array.length)];
const randomDate = (start, end) => new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));

// Data collections
const collections = {
  roles: [],
  departments: [],
  users: [],
  thesis_statuses: [],
  theses: [],
  supervisor_assignments: [],
  submissions: [],
  reviews: [],
  defense_schedules: [],
  defense_scores: [],
  event_logs: [],
  archived_theses: [],
  archived_submissions: [],
  archived_reviews: []
};

// Generate Roles
const generateRoles = () => {
  const roles = [
    {
      _id: new ObjectId(),
      name: 'Sinh viên',
      description: 'Sinh viên thực hiện luận văn',
      permissions: ['submit_thesis', 'view_own_thesis', 'upload_submission'],
      created_at: new Date('2023-01-01')
    },
    {
      _id: new ObjectId(),
      name: 'Giảng viên',
      description: 'Giảng viên hướng dẫn và chấm điểm',
      permissions: ['review_thesis', 'score_thesis', 'view_all_thesis', 'assign_student'],
      created_at: new Date('2023-01-01')
    },
    {
      _id: new ObjectId(),
      name: 'Admin',
      description: 'Quản trị viên hệ thống',
      permissions: ['*'],
      created_at: new Date('2023-01-01')
    }
  ];
  collections.roles = roles;
  console.log(chalk.green('✓ Generated roles'));
};

// Generate Departments
const generateDepartments = () => {
  const depts = [
    { code: 'CNTT', name: 'Công nghệ thông tin', description: 'Khoa Công nghệ thông tin' },
    { code: 'KHMT', name: 'Khoa học máy tính', description: 'Khoa Khoa học máy tính' },
    { code: 'KTPM', name: 'Kỹ thuật phần mềm', description: 'Khoa Kỹ thuật phần mềm' },
    { code: 'HTTT', name: 'Hệ thống thông tin', description: 'Khoa Hệ thống thông tin' },
    { code: 'MMTVTT', name: 'Mạng máy tính và Truyền thông', description: 'Khoa Mạng máy tính và Truyền thông' }
  ];

  collections.departments = depts.map(dept => ({
    _id: new ObjectId(),
    code: dept.code,
    name: dept.name,
    description: dept.description,
    head_id: null, // Will be assigned later
    contact_email: `${dept.code.toLowerCase()}@university.edu.vn`,
    contact_phone: faker.phone.phoneNumber('0## ### ####'),
    is_active: true,
    created_at: new Date('2023-01-01')
  }));
  console.log(chalk.green('✓ Generated departments'));
};

// Generate Thesis Statuses
const generateThesisStatuses = () => {
  const statuses = [
    { name: 'Chờ duyệt', description: 'Đề tài đang chờ phê duyệt', color: '#FFA500', order: 1 },
    { name: 'Đang thực hiện', description: 'Đề tài đang được thực hiện', color: '#0080FF', order: 2 },
    { name: 'Chờ bảo vệ', description: 'Đã nộp bài, chờ bảo vệ', color: '#800080', order: 3 },
    { name: 'Hoàn thành', description: 'Đã bảo vệ thành công', color: '#008000', order: 4 },
    { name: 'Không đạt', description: 'Không đạt yêu cầu', color: '#FF0000', order: 5 }
  ];

  collections.thesis_statuses = statuses.map(status => ({
    _id: new ObjectId(),
    ...status,
    is_active: true
  }));
  console.log(chalk.green('✓ Generated thesis statuses'));
};

// Generate Users
const generateUsers = () => {
  const studentRole = collections.roles.find(r => r.name === 'Sinh viên');
  const teacherRole = collections.roles.find(r => r.name === 'Giảng viên');
  const adminRole = collections.roles.find(r => r.name === 'Admin');
  const majors = ['CNTT', 'KHMT', 'KTPM', 'HTTT', 'MMTVTT'];

  // Generate teachers (20% of users)
  const teacherCount = Math.floor(counts.users * 0.2);
  const studentCount = counts.users - teacherCount - 1; // -1 for admin

  // Admin
  collections.users.push({
    _id: new ObjectId(),
    code: 'ADMIN001',
    full_name: 'Nguyễn Văn Admin',
    email: 'admin@university.edu.vn',
    phone: faker.phone.phoneNumber('0## ### ####'),
    role_id: adminRole._id,
    major: null,
    status: 'active',
    created_at: new Date('2023-01-01'),
    updated_at: new Date()
  });

  // Teachers
  for (let i = 0; i < teacherCount; i++) {
    collections.users.push({
      _id: new ObjectId(),
      code: `GV${String(i + 1).padStart(3, '0')}`,
      full_name: faker.name.findName(),
      email: `gv${String(i + 1).padStart(3, '0')}@university.edu.vn`,
      phone: faker.phone.phoneNumber('0## ### ####'),
      role_id: teacherRole._id,
      major: randomElement(majors),
      status: 'active',
      created_at: randomDate(new Date('2023-01-01'), new Date('2023-06-01')),
      updated_at: new Date()
    });
  }

  // Students
  for (let i = 0; i < studentCount; i++) {
    collections.users.push({
      _id: new ObjectId(),
      code: `SV${String(i + 1).padStart(5, '0')}`,
      full_name: faker.name.findName(),
      email: `sv${String(i + 1).padStart(5, '0')}@student.edu.vn`,
      phone: faker.phone.phoneNumber('0## ### ####'),
      role_id: studentRole._id,
      major: randomElement(majors),
      status: randomElement(['active', 'active', 'active', 'inactive']),
      created_at: randomDate(new Date('2023-01-01'), new Date('2023-09-01')),
      updated_at: new Date()
    });
  }

  // Assign department heads
  const teachers = collections.users.filter(u => u.code.startsWith('GV'));
  collections.departments.forEach((dept, index) => {
    if (index < teachers.length) {
      dept.head_id = teachers[index]._id;
    }
  });

  console.log(chalk.green('✓ Generated users'));
};

// Generate Theses
const generateTheses = () => {
  const students = collections.users.filter(u => u.code.startsWith('SV') && u.status === 'active');
  const teachers = collections.users.filter(u => u.code.startsWith('GV'));
  const statuses = collections.thesis_statuses;
  const academicYears = ['2023-2024', '2024-2025'];
  const semesters = [1, 2];

  const thesisTopics = [
    'Xây dựng hệ thống quản lý',
    'Phát triển ứng dụng di động',
    'Nghiên cứu thuật toán',
    'Tối ưu hóa hiệu năng',
    'Phân tích dữ liệu',
    'Ứng dụng AI/ML trong',
    'Hệ thống IoT cho',
    'Blockchain trong',
    'Bảo mật thông tin cho',
    'Cloud computing cho'
  ];

  const domains = [
    'giáo dục', 'y tế', 'thương mại điện tử', 'logistics',
    'tài chính', 'nông nghiệp', 'du lịch', 'bất động sản',
    'giao thông', 'năng lượng'
  ];

  for (let i = 0; i < counts.theses && i < students.length; i++) {
    const student = students[i];
    const supervisor = randomElement(teachers);
    const status = randomElement(statuses);
    
    collections.theses.push({
      _id: new ObjectId(),
      title: `${randomElement(thesisTopics)} ${randomElement(domains)}`,
      major: student.major,
      description: faker.lorem.paragraphs(2),
      file_url: `/uploads/theses/${new ObjectId()}.pdf`,
      status_id: status._id,
      supervisor_id: supervisor._id,
      student_id: student._id,
      academic_year: randomElement(academicYears),
      semester: randomElement(semesters),
      created_at: randomDate(new Date('2023-09-01'), new Date('2024-02-01')),
      updated_at: new Date()
    });
  }
  console.log(chalk.green('✓ Generated theses'));
};

// Generate Supervisor Assignments
const generateSupervisorAssignments = () => {
  const teachers = collections.users.filter(u => u.code.startsWith('GV'));
  
  collections.theses.forEach(thesis => {
    // Primary supervisor is already assigned
    collections.supervisor_assignments.push({
      _id: new ObjectId(),
      thesis_id: thesis._id,
      supervisor_id: thesis.supervisor_id,
      role: 'primary',
      assigned_at: thesis.created_at,
      assigned_by: collections.users.find(u => u.code === 'ADMIN001')._id,
      is_active: true
    });

    // 30% chance of having a co-supervisor
    if (Math.random() < 0.3) {
      const coSupervisor = randomElement(teachers.filter(t => t._id !== thesis.supervisor_id));
      collections.supervisor_assignments.push({
        _id: new ObjectId(),
        thesis_id: thesis._id,
        supervisor_id: coSupervisor._id,
        role: 'co_supervisor',
        assigned_at: new Date(thesis.created_at.getTime() + 7 * 24 * 60 * 60 * 1000),
        assigned_by: collections.users.find(u => u.code === 'ADMIN001')._id,
        is_active: true
      });
    }
  });
  console.log(chalk.green('✓ Generated supervisor assignments'));
};

// Generate Submissions
const generateSubmissions = () => {
  const submissionTypes = ['midterm', 'final', 'revision'];
  
  collections.theses.forEach(thesis => {
    const submissionCount = Math.floor(Math.random() * 3) + 1;
    let lastSubmissionDate = thesis.created_at;

    for (let i = 0; i < submissionCount && collections.submissions.length < counts.submissions; i++) {
      const submissionDate = randomDate(
        new Date(lastSubmissionDate.getTime() + 30 * 24 * 60 * 60 * 1000),
        new Date(lastSubmissionDate.getTime() + 90 * 24 * 60 * 60 * 1000)
      );

      collections.submissions.push({
        _id: new ObjectId(),
        thesis_id: thesis._id,
        type: i === 0 ? 'midterm' : (i === submissionCount - 1 ? 'final' : 'revision'),
        version: i + 1,
        file_url: `/uploads/submissions/${new ObjectId()}.pdf`,
        file_size: Math.floor(Math.random() * 5000000) + 500000,
        file_name: `${thesis.student_id}_v${i + 1}.pdf`,
        submitted_by: thesis.student_id,
        submitted_at: submissionDate,
        notes: faker.lorem.sentence(),
        status: randomElement(['submitted', 'reviewing', 'approved'])
      });

      lastSubmissionDate = submissionDate;
    }
  });
  console.log(chalk.green('✓ Generated submissions'));
};

// Generate Reviews
const generateReviews = () => {
  const teachers = collections.users.filter(u => u.code.startsWith('GV'));
  const reviewStatuses = ['Pass', 'Fail', 'Revision Required'];

  collections.submissions.forEach((submission, index) => {
    if (index >= counts.reviews) return;

    const thesis = collections.theses.find(t => t._id.equals(submission.thesis_id));
    const reviewers = [thesis.supervisor_id];
    
    // Add co-supervisor if exists
    const coSupervisor = collections.supervisor_assignments.find(
      sa => sa.thesis_id.equals(thesis._id) && sa.role === 'co_supervisor'
    );
    if (coSupervisor) {
      reviewers.push(coSupervisor.supervisor_id);
    }

    reviewers.forEach(reviewerId => {
      collections.reviews.push({
        _id: new ObjectId(),
        submission_id: submission._id,
        reviewer_id: reviewerId,
        score: Math.random() * 3 + 7,
        comment: faker.lorem.sentences(2),
        detailed_feedback: faker.lorem.paragraphs(2),
        status: randomElement(reviewStatuses),
        reviewed_at: randomDate(
          submission.submitted_at,
          new Date(submission.submitted_at.getTime() + 14 * 24 * 60 * 60 * 1000)
        ),
        created_at: submission.submitted_at
      });
    });
  });
  console.log(chalk.green('✓ Generated reviews'));
};

// Generate Defense Schedules
const generateDefenseSchedules = () => {
  const completedTheses = collections.theses
    .filter(t => {
      const status = collections.thesis_statuses.find(s => s._id.equals(t.status_id));
      return status.name === 'Chờ bảo vệ' || status.name === 'Hoàn thành';
    })
    .slice(0, counts.defenses);

  const rooms = ['A101', 'A102', 'B201', 'B202', 'C301'];
  const buildings = ['Tòa A', 'Tòa B', 'Tòa C'];
  const times = ['08:00', '09:00', '10:00', '14:00', '15:00', '16:00'];

  completedTheses.forEach(thesis => {
    const defenseDate = randomDate(new Date('2024-06-01'), new Date('2024-07-31'));
    
    collections.defense_schedules.push({
      _id: new ObjectId(),
      thesis_id: thesis._id,
      defense_date: defenseDate,
      defense_time: randomElement(times),
      duration_minutes: randomElement([30, 45, 60]),
      location: randomElement(rooms),
      building: randomElement(buildings),
      university: 'Đại học Công nghệ Thông tin',
      status: randomElement(['scheduled', 'completed']),
      notes: faker.lorem.sentence(),
      created_at: new Date(defenseDate.getTime() - 30 * 24 * 60 * 60 * 1000),
      updated_at: defenseDate
    });
  });
  console.log(chalk.green('✓ Generated defense schedules'));
};

// Generate Defense Scores
const generateDefenseScores = () => {
  const completedDefenses = collections.defense_schedules.filter(d => d.status === 'completed');
  const teachers = collections.users.filter(u => u.code.startsWith('GV'));
  const criteria = ['presentation', 'content', 'qa_session'];

  completedDefenses.forEach(defense => {
    const thesis = collections.theses.find(t => t._id.equals(defense.thesis_id));
    const committee = [thesis.supervisor_id];
    
    // Add 2-3 more teachers to committee
    const additionalTeachers = teachers
      .filter(t => !t._id.equals(thesis.supervisor_id))
      .sort(() => 0.5 - Math.random())
      .slice(0, Math.floor(Math.random() * 2) + 2);
    
    committee.push(...additionalTeachers.map(t => t._id));

    committee.forEach(scorerId => {
      criteria.forEach(criterion => {
        collections.defense_scores.push({
          _id: new ObjectId(),
          defense_schedule_id: defense._id,
          scorer_id: scorerId,
          score: Math.random() * 2 + 8,
          criteria: criterion,
          comment: faker.lorem.sentences(2),
          scored_at: defense.defense_date
        });
      });
    });
  });
  console.log(chalk.green('✓ Generated defense scores'));
};

// Generate Event Logs
const generateEventLogs = () => {
  const actions = [
    'submit_thesis', 'review_submission', 'schedule_defense',
    'update_thesis', 'assign_supervisor', 'submit_revision'
  ];
  const entityTypes = ['thesis', 'submission', 'review', 'defense'];

  // Generate logs for various activities
  const logCount = Math.min(counts.users * 5, 500);
  
  for (let i = 0; i < logCount; i++) {
    const user = randomElement(collections.users);
    const action = randomElement(actions);
    const entityType = randomElement(entityTypes);
    
    collections.event_logs.push({
      _id: new ObjectId(),
      user_id: user._id,
      action: action,
      entity_type: entityType,
      entity_id: new ObjectId(),
      details: {
        ip: faker.internet.ip(),
        browser: faker.internet.userAgent(),
        additional_info: faker.lorem.sentence()
      },
      ip_address: faker.internet.ip(),
      user_agent: faker.internet.userAgent(),
      timestamp: randomDate(new Date('2023-09-01'), new Date())
    });
  }
  console.log(chalk.green('✓ Generated event logs'));
};

// Generate Archived Data
const generateArchivedData = () => {
  const completedTheses = collections.theses
    .filter(t => {
      const status = collections.thesis_statuses.find(s => s._id.equals(t.status_id));
      return status.name === 'Hoàn thành';
    })
    .slice(0, counts.archived);

  completedTheses.forEach(thesis => {
    const student = collections.users.find(u => u._id.equals(thesis.student_id));
    const supervisor = collections.users.find(u => u._id.equals(thesis.supervisor_id));
    const defenseSchedule = collections.defense_schedules.find(d => d.thesis_id.equals(thesis._id));
    
    const archivedThesis = {
      _id: new ObjectId(),
      original_thesis_id: thesis._id,
      title: thesis.title,
      major: thesis.major,
      description: thesis.description,
      final_file_url: thesis.file_url,
      final_score: Math.random() * 2 + 8,
      graduation_year: 2024,
      supervisor_info: {
        id: supervisor._id,
        name: supervisor.full_name,
        email: supervisor.email
      },
      student_info: {
        id: student._id,
        code: student.code,
        name: student.full_name,
        email: student.email
      },
      archived_at: new Date(),
      archived_by: collections.users.find(u => u.code === 'ADMIN001')._id
    };
    collections.archived_theses.push(archivedThesis);

    // Archive submissions
    const submissions = collections.submissions.filter(s => s.thesis_id.equals(thesis._id));
    submissions.forEach(submission => {
      const archivedSubmission = {
        _id: new ObjectId(),
        archived_thesis_id: archivedThesis._id,
        original_submission_id: submission._id,
        type: submission.type,
        file_url: submission.file_url,
        submitted_at: submission.submitted_at,
        archived_at: new Date()
      };
      collections.archived_submissions.push(archivedSubmission);

      // Archive reviews
      const reviews = collections.reviews.filter(r => r.submission_id.equals(submission._id));
      reviews.forEach(review => {
        const reviewer = collections.users.find(u => u._id.equals(review.reviewer_id));
        collections.archived_reviews.push({
          _id: new ObjectId(),
          archived_submission_id: archivedSubmission._id,
          original_review_id: review._id,
          reviewer_info: {
            id: reviewer._id,
            name: reviewer.full_name,
            email: reviewer.email
          },
          score: review.score,
          comment: review.comment,
          status: review.status,
          reviewed_at: review.reviewed_at,
          archived_at: new Date()
        });
      });
    });
  });
  console.log(chalk.green('✓ Generated archived data'));
};

// Main generation function
const generateAllData = () => {
  console.log(chalk.yellow('\nStarting data generation...'));
  
  generateRoles();
  generateDepartments();
  generateThesisStatuses();
  generateUsers();
  generateTheses();
  generateSupervisorAssignments();
  generateSubmissions();
  generateReviews();
  generateDefenseSchedules();
  generateDefenseScores();
  generateEventLogs();
  generateArchivedData();

  // Save all collections to JSON files
  const dataDir = path.join(__dirname, '..', 'data');
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }

  Object.entries(collections).forEach(([name, data]) => {
    const filePath = path.join(dataDir, `${name}.json`);
    fs.writeFileSync(filePath, JSON.stringify(data, null, 2));
    console.log(chalk.gray(`  Saved ${data.length} ${name} to ${filePath}`));
  });

  console.log(chalk.green('\n✓ Data generation completed!'));
  console.log(chalk.blue(`\nGenerated files in: ${dataDir}`));
};

// Run generation
generateAllData();