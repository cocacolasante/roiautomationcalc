export function calculateROI(answers, config) {
  if (!config || !config.questions) return null;

  const hourlyRate = config.defaultHourlyRate || 35;
  const weeksPerYear = 52;
  const globalAutoPct = config.automationSavingPct || 0.75;

  let totalHoursPerWeek = 0;
  const lineItems = [];

  config.questions.forEach((question) => {
    if (!question.isActive) return;
    const answer = answers[question.key];
    if (answer === undefined || answer === null || answer === '') return;

    const hours = extractHours(question, answer);
    if (!hours || hours <= 0) return;

    const autoPct = question.automationPct || globalAutoPct;
    const annualCost = hours * weeksPerYear * hourlyRate;
    const recoverableHours = hours * autoPct;
    const recoverableCost = recoverableHours * weeksPerYear * hourlyRate;

    totalHoursPerWeek += hours;
    lineItems.push({
      questionKey: question.key,
      description: question.question,
      category: question.category,
      hoursPerWeek: hours,
      annualCost,
      recoverableHours,
      recoverableCost,
    });
  });

  const totalRecoverableHoursPerYear = totalHoursPerWeek * globalAutoPct * weeksPerYear;
  const totalRecoverableCostPerYear = totalRecoverableHoursPerYear * hourlyRate;
  const totalCostPerYear = totalHoursPerWeek * weeksPerYear * hourlyRate;

  const categoryMap = {};
  lineItems.forEach((item) => {
    if (!categoryMap[item.category]) {
      categoryMap[item.category] = { category: item.category, hoursPerWeek: 0, annualCost: 0, recoverableCost: 0 };
    }
    categoryMap[item.category].hoursPerWeek += item.hoursPerWeek;
    categoryMap[item.category].annualCost += item.annualCost;
    categoryMap[item.category].recoverableCost += item.recoverableCost;
  });

  const categories = Object.values(categoryMap).sort((a, b) => b.recoverableCost - a.recoverableCost);

  return {
    totalHoursPerWeek: Math.round(totalHoursPerWeek * 10) / 10,
    totalHoursPerYear: Math.round(totalHoursPerWeek * weeksPerYear),
    totalCostPerYear: Math.round(totalCostPerYear),
    recoverableHoursPerYear: Math.round(totalRecoverableHoursPerYear),
    recoverableCostPerYear: Math.round(totalRecoverableCostPerYear),
    categories,
    lineItems,
    hourlyRate,
  };
}

function extractHours(question, answer) {
  switch (question.inputType) {
    case 'hours_slider':
      return parseFloat(answer) || 0;
    case 'number':
    case 'people_count': {
      const count = parseFloat(answer) || 0;
      return count * (question.hoursPerUnit || 1);
    }
    case 'multiple_choice': {
      const opt = (question.options || []).find((o) => o.value === answer);
      return opt ? opt.hours || 0 : 0;
    }
    case 'yes_no':
      return answer === 'yes' ? question.hoursPerUnit || 0 : 0;
    default:
      return 0;
  }
}

export function formatMoney(amount) {
  if (amount >= 1000000) return `$${(amount / 1000000).toFixed(1)}M`;
  if (amount >= 1000) return `$${Math.round(amount / 1000)}K`;
  return `$${Math.round(amount)}`;
}
