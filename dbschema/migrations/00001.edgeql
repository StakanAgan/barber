CREATE MIGRATION m17afwt5kee6rq33v7knsn4xormlvxh4kx2il55eu4jerclz2gd6mq
    ONTO initial
{
  CREATE SCALAR TYPE default::ServiceType EXTENDING enum<Hair, Beard, HairBeard>;
  CREATE TYPE default::Barber {
      CREATE REQUIRED PROPERTY availableTypes -> array<default::ServiceType>;
      CREATE REQUIRED PROPERTY fullName -> std::str;
      CREATE REQUIRED PROPERTY phone -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY telegramId -> std::int64 {
          CREATE CONSTRAINT std::exclusive;
      };
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
      CREATE MULTI LINK shifts -> default::BarberShift;
  };
  CREATE SCALAR TYPE default::VisitStatus EXTENDING enum<Created, InProcess, Done, Canceled>;
  CREATE TYPE default::Visit {
      CREATE REQUIRED LINK barberShift -> default::BarberShift;
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
      CREATE REQUIRED PROPERTY serviceType -> default::ServiceType;
      CREATE PROPERTY status -> default::VisitStatus;
  };
  ALTER TYPE default::BarberShift {
      CREATE MULTI LINK visits -> default::Visit;
  };
  CREATE TYPE default::Customer {
      CREATE MULTI LINK visits -> default::Visit;
      CREATE REQUIRED PROPERTY fullName -> std::str;
      CREATE REQUIRED PROPERTY phone -> std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY telegramId -> std::int64 {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE default::Visit {
      CREATE REQUIRED LINK customer -> default::Customer;
  };
};
