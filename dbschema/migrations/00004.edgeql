CREATE MIGRATION m1g5m376qkt76sjbggqisyjqwl53xkv4ou6utu4utf7wq7u66jflaa
    ONTO m17qjj2ojt2ugia2jwafxzm7vsyjrjfdz4wgv5voqghyfzgwvq7eeq
{
  ALTER TYPE default::Visit {
      CREATE CONSTRAINT std::exclusive ON ((.plannedFrom, .barberShift));
  };
};
