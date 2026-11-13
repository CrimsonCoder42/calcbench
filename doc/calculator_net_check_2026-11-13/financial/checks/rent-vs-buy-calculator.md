# rent-vs-buy-calculator

FIELDS (form `calform`): chomeprice (500,000), cdownpay (20 %), cinterest (7.263 %), cloanterm (30), cbuyclosing (2 %),
cpropertytax (1.5 %/yr) + increase (3 %), chomeinsurance (2,500/yr), choa (0), cmaintenance (1.5 %/yr),
cvalueincrease (3 %/yr), ccostinsuranceincrease (3 %), csellclosing (7 %), crental (3,000/mo) + increase (3 %),
crentinsurance (15/mo), cdeposit (3,000), cupfront (100), cinvestreturn (5 %), cfedtax (25 %), cstatetax (0),
cfilestatus (Single|MarriedJoint|MarriedSeparately|HeadofHousehold|Widow).

CLASSIFY: SKIPPED. The filing-status select means the model applies the federal standard deduction to decide the
mortgage-interest deduction, a figure that changes every year. Without it the page is a multi-year scenario built
from pieces we already have (FIN-INT-018 payment, FIN-LOAN-002 balance, FIN-INT-001 growth, FIN-INT-038 NPV).
No submissions.
