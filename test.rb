class Test
  def this_test
    puts "hello #{test2('from 1')}"
  end

  def test2(a)
    a
  end
end

Test.new.this_test
