// Discount Engine — Strategy + Decorator reference solution (JavaScript).

class CartItem {
  constructor(name, price, quantity, category = "") {
    this.name = name;
    this.price = price;
    this.quantity = quantity;
    this.category = category;
  }
}

function cartTotal(cart) {
  let total = 0.0;
  for (const i of cart) total += i.price * i.quantity;
  return total;
}

class Discount {
  apply(cart) { throw new Error('not implemented'); }
}

class PercentageDiscount extends Discount {
  constructor(pct) {
    super();
    this.pct = pct;
  }
  apply(cart) {
    const total = cartTotal(cart);
    return total * (1.0 - this.pct / 100.0);
  }
}

class FlatDiscount extends Discount {
  constructor(amount) {
    super();
    this.amount = amount;
  }
  apply(cart) {
    const total = cartTotal(cart);
    return Math.max(0.0, total - this.amount);
  }
}

class BuyXGetYDiscount extends Discount {
  constructor(buy, free) {
    super();
    this.buy = buy;
    this.free = free;
  }
  apply(cart) {
    const group = this.buy + this.free;
    let total = 0.0;
    for (const it of cart) {
      const groups = Math.floor(it.quantity / group);
      const remainder = it.quantity % group;
      const paid = groups * this.buy + Math.min(remainder, this.buy);
      total += paid * it.price;
    }
    return total;
  }
}

class StackedDiscount extends Discount {
  constructor(discounts) {
    super();
    this.discounts = discounts;
  }
  apply(cart) {
    let current = cartTotal(cart);
    for (const d of this.discounts) {
      const temp = [new CartItem("subtotal", current, 1, "")];
      current = d.apply(temp);
    }
    return current;
  }
}

function apply_percentage_discount(cart, percentage) {
  return new PercentageDiscount(percentage).apply(cart);
}

function apply_flat_discount(cart, amount) {
  return new FlatDiscount(amount).apply(cart);
}

function apply_bogo(cart, buy, free) {
  return new BuyXGetYDiscount(buy, free).apply(cart);
}

function apply_percentage_with_eligibility(cart, percentage, minCartValue,
    requireFirstTimeUser, isFirstTimeUser, eligibleCategory) {
  const raw = cartTotal(cart);
  if (raw < minCartValue) return raw;
  if (requireFirstTimeUser && !isFirstTimeUser) return raw;
  if (eligibleCategory) {
    const eligible = cart.filter(i => i.category === eligibleCategory);
    let nonEligible = 0.0;
    for (const i of cart) {
      if (i.category !== eligibleCategory) nonEligible += i.price * i.quantity;
    }
    return new PercentageDiscount(percentage).apply(eligible) + nonEligible;
  }
  return new PercentageDiscount(percentage).apply(cart);
}

function reset_service() {}

module.exports = {
  CartItem,
  Discount,
  PercentageDiscount,
  FlatDiscount,
  BuyXGetYDiscount,
  StackedDiscount,
  reset_service,
  apply_percentage_discount,
  apply_flat_discount,
  apply_bogo,
  apply_percentage_with_eligibility,
};
