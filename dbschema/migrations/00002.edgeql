CREATE MIGRATION m1tzuexzo57a4ginhadlpqnai7po376owcfe4wac3pd34tyfo7mlyq
    ONTO m17afwt5kee6rq33v7knsn4xormlvxh4kx2il55eu4jerclz2gd6mq
{
  ALTER TYPE default::Barber {
      CREATE REQUIRED PROPERTY timeZoneOffset -> std::int64 {
          SET default := 3;
      };
  };
  ALTER TYPE default::Customer {
      CREATE REQUIRED PROPERTY timeZoneOffset -> std::int64 {
          SET default := 3;
      };
  };
};
