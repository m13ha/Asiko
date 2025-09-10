/**
 * Basic usage example for the Appointment Master API Client
 * 
 * This example demonstrates:
 * - User registration and login
 * - Creating appointments
 * - Booking appointments
 * - Managing bookings
 */

const { AppointmentMasterClient } = require('../dist/client');

async function basicExample() {
  // Initialize the client
  const client = new AppointmentMasterClient({
    baseUrl: 'http://localhost:8888'
  });

  try {
    console.log('🚀 Starting Appointment Master API Demo...\n');

    // 1. Register a new user
    console.log('1. Registering new user...');
    const newUser = await client.register({
      name: 'John Doe',
      email: 'john.doe@example.com',
      password: 'securepassword123',
      phoneNumber: '+1234567890'
    });
    console.log('✅ User registered:', newUser.name);

    // 2. Login
    console.log('\n2. Logging in...');
    const loginResponse = await client.login('john.doe@example.com', 'securepassword123');
    console.log('✅ Login successful! Token received.');

    // 3. Create an appointment
    console.log('\n3. Creating appointment...');
    const appointment = await client.createAppointment({
      title: 'Consultation Session',
      startTime: new Date('2025-01-15T09:00:00Z'),
      endTime: new Date('2025-01-15T17:00:00Z'),
      startDate: new Date('2025-01-15T00:00:00Z'),
      endDate: new Date('2025-01-20T00:00:00Z'),
      bookingDuration: 60,
      type: 'single',
      maxAttendees: 1,
      description: 'Professional consultation'
    });
    console.log('✅ Appointment created:', appointment.appCode);

    // 4. Get available slots
    console.log('\n4. Getting available slots...');
    const slots = await client.getAvailableSlots(appointment.appCode);
    console.log('✅ Available slots found:', slots.items?.length || 0);

    console.log('\n🎉 Demo completed successfully!');

  } catch (error) {
    console.error('❌ Error:', error.message);
  }
}

// Run the example
if (require.main === module) {
  basicExample();
}

module.exports = { basicExample };