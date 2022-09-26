CREATE MIGRATION m1j2gugfraf2udhjtka3fpq2bt55cuyyvx2dobxbzsdfwilirvlhvq
    ONTO m1pmw7e4keqvq4bxdbrunqjzgn62uxjybc7gk4c34lnv7iweakp5wq
{
  ALTER TYPE default::Service {
      DROP PROPERTY type;
  };
  DROP SCALAR TYPE default::ServiceType;
};
