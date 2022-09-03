CREATE MIGRATION m1pmw7e4keqvq4bxdbrunqjzgn62uxjybc7gk4c34lnv7iweakp5wq
    ONTO initial
{
  CREATE TYPE default::Barber {
      CREATE REQUIRED PROPERTY fullName -> std::str;
      CREATE REQUIRED PROPERTY phone -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY telegramId -> std::int64 {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY timeZoneOffset -> std::int64 {
          SET default := 3;
      };
  };
  CREATE SCALAR TYPE default::ServiceType EXTENDING enum<Hair, Beard, HairBeard>;
  CREATE TYPE default::Service {
      CREATE REQUIRED LINK barber -> default::Barber;
      CREATE REQUIRED PROPERTY duration -> std::duration;
      CREATE REQUIRED PROPERTY price -> std::int64 {
          CREATE CONSTRAINT std::min_value(0);
      };
      CREATE REQUIRED PROPERTY title -> std::str;
      CREATE REQUIRED PROPERTY type -> default::ServiceType;
  };
  ALTER TYPE default::Barber {
      CREATE MULTI LINK services := (.<barber[IS default::Service]);
  };
  CREATE SCALAR TYPE default::ShiftStatus EXTENDING enum<Planned, Work, Finished>;
  CREATE TYPE default::BarberShift {
      CREATE REQUIRED LINK barber -> default::Barber;
      CREATE PROPERTY actualFrom -> std::datetime;
      CREATE PROPERTY actualTo -> std::datetime;
      CREATE CONSTRAINT std::expression ON ((.actualFrom < .actualTo));
      CREATE REQUIRED PROPERTY plannedFrom -> std::datetime;
      CREATE REQUIRED PROPERTY plannedTo -> std::datetime;
      CREATE CONSTRAINT std::expression ON ((.plannedFrom < .plannedTo));
      CREATE REQUIRED PROPERTY status -> default::ShiftStatus;
  };
  ALTER TYPE default::Barber {
      CREATE MULTI LINK shifts := (.<barber[IS default::BarberShift]);
  };
  CREATE SCALAR TYPE default::VisitStatus EXTENDING enum<Created, InProcess, Done, Canceled>;
  CREATE TYPE default::Visit {
      CREATE REQUIRED LINK barberShift -> default::BarberShift;
      CREATE REQUIRED LINK service -> default::Service;
      CREATE PROPERTY actualFrom -> std::datetime;
      CREATE PROPERTY actualTo -> std::datetime;
      CREATE CONSTRAINT std::expression ON ((.actualFrom < .actualTo));
      CREATE PROPERTY discountPrice -> std::int64 {
          CREATE CONSTRAINT std::min_value(0);
      };
      CREATE PROPERTY price -> std::int64 {
          CREATE CONSTRAINT std::min_value(0);
      };
      CREATE CONSTRAINT std::expression ON ((.discountPrice <= .price));
      CREATE REQUIRED PROPERTY plannedFrom -> std::datetime;
      CREATE REQUIRED PROPERTY plannedTo -> std::datetime;
      CREATE CONSTRAINT std::expression ON ((.plannedFrom < .plannedTo));
      CREATE REQUIRED PROPERTY deleted -> std::bool {
          SET default := false;
      };
      CREATE PROPERTY totalPrice := ((.price - .discountPrice));
      CREATE PROPERTY status -> default::VisitStatus;
  };
  ALTER TYPE default::BarberShift {
      CREATE MULTI LINK visits := (.<barberShift[IS default::Visit]);
  };
  CREATE TYPE default::Customer {
      CREATE REQUIRED PROPERTY fullName -> std::str;
      CREATE REQUIRED PROPERTY phone -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY telegramId -> std::int64 {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY timeZoneOffset -> std::int64 {
          SET default := 3;
      };
  };
  ALTER TYPE default::Visit {
      CREATE REQUIRED LINK customer -> default::Customer;
  };
  ALTER TYPE default::Customer {
      CREATE MULTI LINK visits := (.<customer[IS default::Visit]);
  };
};
