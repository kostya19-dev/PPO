const databaseName = process.env.MONGODB_DATABASE || "vitaline_clinic";
const clinic = db.getSiblingDB(databaseName);

clinic.patients.deleteMany({});
clinic.appointments.deleteMany({});
clinic.schedule_slots.deleteMany({});
clinic.visits.deleteMany({});
clinic.emergency_requests.deleteMany({});
clinic.reschedule_offers.deleteMany({});

clinic.patients.createIndex({ email: 1 }, { unique: true });
clinic.appointments.createIndex(
  { doctor_id: 1, appointment_date: 1 },
  {
    unique: true,
    partialFilterExpression: { status: { $in: ["CREATED", "COMPLETED"] } },
  }
);
clinic.appointments.createIndex({ patient_id: 1, appointment_date: 1 });
clinic.schedule_slots.createIndex({ doctor_id: 1, slot_date: 1, start_time: 1 });
clinic.visits.createIndex({ appointment_id: 1 }, { unique: true });
clinic.emergency_requests.createIndex({ status: 1, requested_date: 1 });
clinic.reschedule_offers.createIndex({ emergency_request_id: 1, status: 1 });

const tomorrow = new Date();
tomorrow.setUTCDate(tomorrow.getUTCDate() + 1);
tomorrow.setUTCHours(0, 0, 0, 0);

function at(hour, minute) {
  const value = new Date(tomorrow);
  value.setUTCHours(hour, minute, 0, 0);
  return value;
}

clinic.patients.insertMany([
  {
    _id: "11111111-1111-1111-1111-111111111111",
    full_name: "Иванов Иван Иванович",
    phone: "+79001234567",
    email: "ivanov@example.com",
    birth_date: ISODate("1990-05-15T00:00:00Z"),
    created_at: new Date(),
  },
  {
    _id: "22222222-2222-2222-2222-222222222222",
    full_name: "Петров Пётр Петрович",
    phone: "+79007654321",
    email: "petrov@example.com",
    birth_date: ISODate("1985-03-20T00:00:00Z"),
    created_at: new Date(),
  },
  {
    _id: "99999999-9999-9999-9999-999999999999",
    full_name: "Смирнова Анна Викторовна",
    phone: "+79005554433",
    email: "smirnova@example.com",
    birth_date: ISODate("1995-08-11T00:00:00Z"),
    created_at: new Date(),
  },
]);

clinic.appointments.insertMany([
  {
    _id: "33333333-3333-3333-3333-333333333333",
    patient_id: "11111111-1111-1111-1111-111111111111",
    doctor_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    status: "CREATED",
    appointment_date: at(9, 0),
    created_at: new Date(),
  },
  {
    _id: "44444444-4444-4444-4444-444444444444",
    patient_id: "22222222-2222-2222-2222-222222222222",
    doctor_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    status: "COMPLETED",
    appointment_date: at(10, 0),
    created_at: new Date(),
  },
]);

clinic.schedule_slots.insertMany([
  {
    _id: "55555555-5555-5555-5555-555555555555",
    doctor_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    slot_date: tomorrow,
    start_time: at(9, 0),
    end_time: at(9, 30),
    cabinet: "101",
    is_booked: true,
  },
  {
    _id: "66666666-6666-6666-6666-666666666666",
    doctor_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    slot_date: tomorrow,
    start_time: at(9, 30),
    end_time: at(10, 0),
    cabinet: "101",
    is_booked: false,
  },
  {
    _id: "88888888-8888-8888-8888-888888888888",
    doctor_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    slot_date: tomorrow,
    start_time: at(10, 0),
    end_time: at(10, 30),
    cabinet: "202",
    is_booked: true,
  },
]);

clinic.visits.insertOne({
  _id: "77777777-7777-7777-7777-777777777777",
  appointment_id: "44444444-4444-4444-4444-444444444444",
  complaints: "Головная боль, слабость",
  diagnosis: "ОРВИ",
  recommendations: "Покой, питьевой режим, контроль температуры",
  status: "COMPLETED",
  created_at: new Date(),
});

print(`MongoDB migration completed for ${databaseName}`);
