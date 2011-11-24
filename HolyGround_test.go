package main

import "testing"

func TestHolyGroundCalculate(t *testing.T) {
    terrain := NewTerrain(
    "0.....................................................\n" +
    ".%%%%%%%%%%...........................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    ".%....................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................\n" +
    "......................................................")

    expected := 
    "0123456789abcdefghijklmnopqrqponmlkjihgfedcba987654321\n" +
    "1++++++++++cdefghijklmnopqrsrqponmlkjihgfedcba98765432\n" +
    "2+klkjihgfedefghijklmnopqrstsrqponmlkjihgfedcba9876543\n" +
    "3+jklkjihgfefghijklmnopqrstutsrqponmlkjihgfedcba987654\n" +
    "4+ijklkjihgfghijklmnopqrstu+utsrqponmlkjihgfedcba98765\n" +
    "5+hijklkjihghijklmnopqrstu+++utsrqponmlkjihgfedcba9876\n" +
    "6+ghijklkjihijklmnopqrstu+++++utsrqponmlkjihgfedcba987\n" +
    "7+fghijklkjijklmnopqrstu+++++++utsrqponmlkjihgfedcba98\n" +
    "8+efghijklkjklmnopqrstu+++++++++utsrqponmlkjihgfedcba9\n" +
    "9+defghijklklmnopqrstu+++++++++++utsrqponmlkjihgfedcba\n" +
    "abcdefghijklmnopqrstu+++++++++++++utsrqponmlkjihgfedcb\n" +
    "abcdefghijklmnopqrstu+++++++++++++utsrqponmlkjihgfedcb\n" +
    "9abcdefghijklmnopqrstu+++++++++++utsrqponmlkjihgfedcba\n" +
    "89abcdefghijklmnopqrstu+++++++++utsrqponmlkjihgfedcba9\n" +
    "789abcdefghijklmnopqrstu+++++++utsrqponmlkjihgfedcba98\n" +
    "6789abcdefghijklmnopqrstu+++++utsrqponmlkjihgfedcba987\n" +
    "56789abcdefghijklmnopqrstu+++utsrqponmlkjihgfedcba9876\n" +
    "456789abcdefghijklmnopqrstu+utsrqponmlkjihgfedcba98765\n" +
    "3456789abcdefghijklmnopqrstutsrqponmlkjihgfedcba987654\n" +
    "23456789abcdefghijklmnopqrstsrqponmlkjihgfedcba9876543\n" +
    "123456789abcdefghijklmnopqrsrqponmlkjihgfedcba98765432"

    holyGround := NewHolyGround(terrain)

    holyGround.Calculate()

    if holyGround.String() != expected {
        t.Error(holyGround)
    }
}
