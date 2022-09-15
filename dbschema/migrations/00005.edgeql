CREATE MIGRATION m1ncl4meqhiq6mfhuanmhcnaopwiqiv2uyfy4fa6yp2lq4efp2khfa
    ONTO m1g5m376qkt76sjbggqisyjqwl53xkv4ou6utu4utf7wq7u66jflaa
{
  ALTER TYPE default::Visit {
      ALTER PROPERTY discountPrice {
          SET default := 0;
          SET REQUIRED USING (0);
      };
  };
};
