CREATE MIGRATION m17qjj2ojt2ugia2jwafxzm7vsyjrjfdz4wgv5voqghyfzgwvq7eeq
    ONTO m1j2gugfraf2udhjtka3fpq2bt55cuyyvx2dobxbzsdfwilirvlhvq
{
  ALTER TYPE default::Service {
      CREATE REQUIRED PROPERTY deleted -> std::bool {
          SET default := false;
      };
  };
};
