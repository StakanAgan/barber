module default {
  type Customer {
    required property fullName -> str;
    required property phone -> str{
      constraint exclusive;
    };
    required property telegramId -> int64{
      constraint exclusive
    };
    required property timeZoneOffset -> int64{
      default := 3;
    };
    multi link visits := .<customer[is Visit]
  }

  type Barber {
    required property fullName -> str;
    required property phone -> str{
      constraint exclusive;
    };
    required property telegramId -> int64{
      constraint exclusive;
    };
    required property timeZoneOffset -> int64{
      default := 3;
    };
    multi link shifts := .<barber[is BarberShift];
    multi link services := .<barber[is Service] 
  }
  scalar type ShiftStatus extending enum<Planned, Work, Finished>;
  type BarberShift {
    required link barber -> Barber;
    required property status -> ShiftStatus;
    required property plannedFrom -> datetime;
    required property plannedTo -> datetime;
    constraint expression on (
      .plannedFrom < .plannedTo
    );
    property actualFrom -> datetime;
    property actualTo -> datetime;
    constraint expression on (
      .actualFrom < .actualTo
    );
    multi link visits := .<barberShift[is Visit]
  }
  scalar type VisitStatus extending enum<Created, InProcess, Done, Canceled>;
  type Visit {
    required link customer -> Customer;
    required link barberShift -> BarberShift;
    required property plannedFrom -> datetime;
    required property plannedTo -> datetime;
    constraint expression on (
      .plannedFrom < .plannedTo
    );
    property actualFrom -> datetime;
    property actualTo -> datetime;
    constraint expression on (
      .actualFrom < .actualTo
    );
    required link service -> Service;
    property price -> int64{
      constraint min_value(0);  
    };
    property discountPrice -> int64{
      constraint min_value(0);
    };
    constraint expression on (
      .discountPrice <= .price
    );
    property totalPrice := .price - .discountPrice;
    property status -> VisitStatus;
    required property deleted -> bool{
      default := false;
    };
  }
  type Service {
    required link barber -> Barber;
    required property title -> str;
    required property price -> int64{
      constraint min_value(0);
    };
    required property duration -> duration;
    required property deleted -> bool{
      default := false;
    }
  }
}
