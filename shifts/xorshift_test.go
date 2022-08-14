package xor

import (
	"testing"
)

func TestShift64(t *testing.T) {
	shift_a := NewShift64(0)

	if shift_a.seed != 5036377382042008862 {
		t.Fatalf(`seed failed, got %d`, shift_a.seed)
	}

	sa := shift_a.shift();
	if sa != -6399782287330682226 {
		t.Fatalf(`sa failed, got %d`, sa)
	}
	sb := shift_a.shift();
	if sb != 4297237695309840522 {
		t.Fatalf(`sb failed, got %d`, sb)
	}
	sc := shift_a.shift();
	if sc != 1075437695011947220 {
		t.Fatalf(`sc failed, got %d`, sc)
	}
	sd := shift_a.shift();
	if sd != -930821246400571898 {
		t.Fatalf(`sd failed, got %d`, sd)
	}

	ia := shift_a.Int32(2048);
	if ia != 1025 {
		t.Fatalf(`ia failed, got %d`, ia)
	}
	ib := shift_a.Int32(1024);
	if ib != 798 {
		t.Fatalf(`ib failed, got %d`, ib)
	}
	ic := shift_a.Int32(512);
	if ic != 235 {
		t.Fatalf(`ic failed, got %d`, ic)
	}
	id := shift_a.Int32(256);
	if id != 205 {
		t.Fatalf(`se failed, got %d`, id)
	}

	isa := shift_a.Int32(128);
	if isa != 108 {
		t.Fatalf(`isa failed, got %d`, isa)
	}
	isb := shift_a.Int32(64);
	if isb != 42 {
		t.Fatalf(`isb failed, got %d`, isb)
	}
	isc := shift_a.Int32(128);
	if isc != 30 {
		t.Fatalf(`isc failed, got %d`, isc)
	}

	ie := shift_a.Int32(256);
	if ie != 185 {
		t.Fatalf(`ie failed, got %d`, ie)
	}
	ig := shift_a.Int32(512);
	if ig != 237 {
		t.Fatalf(`ig failed, got %d`, ig)
	}
	ih := shift_a.Int32(1024);
	if ih != 974 {
		t.Fatalf(`ih failed, got %d`, ih)
	}
	ii := shift_a.Int32(2048);
	if ii != 1385 {
		t.Fatalf(`ii failed, got %d`, ii)
	}

	se := shift_a.shift();
	if se != -5828336445164884370 {
		t.Fatalf(`se failed, got %d`, se)
	}
	sf := shift_a.shift();
	if sf != 1599167847083165552 {
		t.Fatalf(`sf failed, got %d`, sf)
	}
	sg := shift_a.shift();
	if sg != 6218638069927327200 {
		t.Fatalf(`sg failed, got %d`, sg)
	}
	sh := shift_a.shift();
	if sh != 8232039552211122488 {
		t.Fatalf(`sh failed, got %d`, sh)
	}
}

func TestGetSeed(t *testing.T) {
	a := GetSeed("")
	b := GetSeed("a")
	c := GetSeed("bb")
	d := GetSeed("ccc")
	e := GetSeed("near")

	if a != 5698237097726351552 {
		t.Fatalf(`seed on 'a' failed, got %d`, a)
	}
	if b != 3027204654264679692 {
		t.Fatalf(`seed on 'b' failed, got %d`, b)
	}
	if c != 331832489265128583 {
		t.Fatalf(`seed on 'c' failed, got %d`, c)
	}
	if d != 1883749424214749104 {
		t.Fatalf(`seed on 'd' failed, got %d`, d)
	}
	if e != -4661580130154814320 {
		t.Fatalf(`seed on 'e' failed, got %d`, e)
	}
}
