CREATE MIGRATION m1fsleiwpv7jutmljovcln5sgxw6k3mf7qqnioholgeyimu7j3diga
    ONTO m1ncl4meqhiq6mfhuanmhcnaopwiqiv2uyfy4fa6yp2lq4efp2khfa
{
  ALTER SCALAR TYPE default::ShiftStatus EXTENDING enum<Planned, Work, Finished, Canceled>;
};
