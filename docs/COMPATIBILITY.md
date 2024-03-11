# Cosmium Compatibility with Cosmos DB

## Introduction

Cosmium is designed to emulate the functionality of Cosmos DB, providing developers with a local development environment that closely mimics the behavior of Cosmos DB. While Cosmium aims to be compatible with Cosmos DB, there are certain differences and limitations to be aware of. This document provides an overview of Cosmium's compatibility with Cosmos DB and highlights areas where deviations may occur.

## Supported Features

Cosmium strives to support the core features of Cosmos DB, including:

- REST API
- SQL-like query language
- Document-based data model

## Compatibility Matrix

### Features
| Feature                       | Implemented |
|-------------------------------|-------------|
| Subqueries                    | No          |
| Joins                         | No          |
| Computed properties           | No          |
| Coalesce operators            | No          |
| Bitwise operators             | No          |
| GeoJSON location data         | No          |
| Parameterized queries         | Yes         |
| Stored procedures             | No          |
| Triggers                      | No          |
| User-defined functions (UDFs) | No          |

### Clauses
| Clause       | Implemented |
|--------------|-------------|
| SELECT       | Yes         |
| FROM         | Yes         |
| WHERE        | Yes         |
| ORDER BY     | Yes         |
| GROUP BY     | Yes         |
| OFFSET LIMIT | No          |

### Keywords
| Keyword  | Implemented |
|----------|-------------|
| BETWEEN  | No          |
| DISTINCT | Yes         |
| LIKE     | No          |
| IN       | Yes         |
| TOP      | Yes         |

### Aggregate Functions
| Function | Implemented |
|----------|-------------|
| AVG      | Yes         |
| COUNT    | Yes         |
| MAX      | Yes         |
| MIN      | Yes         |
| SUM      | Yes         |

### Array Functions
| Function       | Implemented |
|----------------|-------------|
| ARRAY_CONCAT   | Yes         |
| ARRAY_CONTAINS | No          |
| ARRAY_LENGTH   | Yes         |
| ARRAY_SLICE    | Yes         |
| CHOOSE         | No          |
| ObjectToArray  | No          |
| SetIntersect   | Yes         |
| SetUnion       | Yes         |

### Conditional Functions
| Function | Implemented |
|----------|-------------|
| IIF      | No          |

### Date and time Functions
| Function                  | Implemented |
|---------------------------|-------------|
| DateTimeAdd               | No          |
| DateTimeBin               | No          |
| DateTimeDiff              | No          |
| DateTimeFromParts         | No          |
| DateTimePart              | No          |
| DateTimeToTicks           | No          |
| DateTimeToTimestamp       | No          |
| GetCurrentDateTime        | No          |
| GetCurrentDateTimeStatic  | No          |
| GetCurrentTicks           | No          |
| GetCurrentTicksStatic     | No          |
| GetCurrentTimestamp       | No          |
| GetCurrentTimestampStatic | No          |
| TicksToDateTime           | No          |
| TimestampToDateTime       | No          |

### Item Functions
| Function   | Implemented |
|------------|-------------|
| DocumentId | No          |

### Mathematical Functions
| Function         | Implemented |
|------------------|-------------|
| ABS              | No          |
| ACOS             | No          |
| ASIN             | No          |
| ATAN             | No          |
| ATN2             | No          |
| CEILING          | No          |
| COS              | No          |
| COT              | No          |
| DEGREES          | No          |
| EXP              | No          |
| FLOOR            | No          |
| IntAdd           | No          |
| IntBitAnd        | No          |
| IntBitLeftShift  | No          |
| IntBitNot        | No          |
| IntBitOr         | No          |
| IntBitRightShift | No          |
| IntBitXor        | No          |
| IntDiv           | No          |
| IntMod           | No          |
| IntMul           | No          |
| IntSub           | No          |
| LOG              | No          |
| LOG10            | No          |
| NumberBin        | No          |
| PI               | No          |
| POWER            | No          |
| RADIANS          | No          |
| RAND             | No          |
| ROUND            | No          |
| SIGN             | No          |
| SIN              | No          |
| SQRT             | No          |
| SQUARE           | No          |
| TAN              | No          |
| TRUNC            | No          |

### Spatial Functions
| Function           | Implemented |
|--------------------|-------------|
| ST_AREA            | No          |
| ST_DISTANCE        | No          |
| ST_WITHIN          | No          |
| ST_INTERSECTS      | No          |
| ST_ISVALID         | No          |
| ST_ISVALIDDETAILED | No          |

### String Functions
| Function        | Implemented |
|-----------------|-------------|
| CONCAT          | Yes         |
| CONTAINS        | Yes         |
| ENDSWITH        | Yes         |
| INDEX_OF        | Yes         |
| LEFT            | Yes         |
| LENGTH          | Yes         |
| LOWER           | Yes         |
| LTRIM           | Yes         |
| REGEXMATCH      | No          |
| REPLACE         | Yes         |
| REPLICATE       | Yes         |
| REVERSE         | Yes         |
| RIGHT           | Yes         |
| RTRIM           | Yes         |
| STARTSWITH      | Yes         |
| STRINGEQUALS    | Yes         |
| StringToArray   | No          |
| StringToBoolean | No          |
| StringToNull    | No          |
| StringToNumber  | No          |
| StringToObject  | No          |
| SUBSTRING       | Yes         |
| ToString        | Yes         |
| TRIM            | Yes         |
| UPPER           | Yes         |

### Type checking Functions
| Function         | Implemented |
|------------------|-------------|
| IS_ARRAY         | Yes         |
| IS_BOOL          | Yes         |
| IS_DEFINED       | Yes         |
| IS_FINITE_NUMBER | Yes         |
| IS_INTEGER       | Yes         |
| IS_NULL          | Yes         |
| IS_NUMBER        | Yes         |
| IS_OBJECT        | Yes         |
| IS_PRIMITIVE     | Yes         |
| IS_STRING        | Yes         |

## Known Differences

While Cosmium aims to replicate the behavior of Cosmos DB as closely as possible, there are certain differences and limitations to be aware of:

1. **Performance**: Cosmium may exhibit different performance characteristics compared to Cosmos DB, especially under heavy load or large datasets.
2. **Consistency Levels**: The consistency model in Cosmium may differ slightly from Cosmos DB.
3. **Features**: Some advanced features or functionalities of Cosmos DB may not be fully supported or available in Cosmium.

## Future Development

Cosmium is actively developed and maintained, with ongoing efforts to improve compatibility with Cosmos DB and enhance its features and capabilities. Future updates may address known differences and limitations, as well as introduce new functionality to bring Cosmium closer to feature parity with Cosmos DB.
