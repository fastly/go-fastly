package fastly

import "testing"

func TestClient_GetBilling(t *testing.T) {
	b, err := testClient.GetBilling(&GetBillingInput{
		Year:  2015,
		Month: 11,
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.InvoiceID == "" {
		t.Errorf("bad invoice_id: %q", b.InvoiceID)
	}
	if b.StartTime == nil {
		t.Errorf("bad start_time: %q", b.StartTime)
	}
	if b.EndTime == nil {
		t.Errorf("bad end_time: %q", b.EndTime)
	}
	if b.Status == nil {
		t.Errorf("bad status: %v", b.Status)
	}
	if b.Total == nil {
		t.Errorf("bad total: %v", b.Total)
	}

	if b.Status.InvoiceID == "" {
		t.Errorf("bad status.invoice_id: %q", b.Status.InvoiceID)
	}
	if b.Status.Status == "" {
		t.Errorf("bad status.status: %q", b.Status.Status)
	}
	if b.Status.SentAt == nil {
		t.Errorf("bad status.sent_at: %q", b.Status.SentAt)
	}

	if b.Total.PlanName == "" {
		t.Errorf("bad total.plan_name: %q", b.Total.PlanName)
	}
	if b.Total.PlanCode == "" {
		t.Errorf("bad total.plan_code: %q", b.Total.PlanCode)
	}
	if b.Total.PlanMinimum == "" {
		t.Errorf("bad total.plan_minimum: %q", b.Total.PlanMinimum)
	}
}
